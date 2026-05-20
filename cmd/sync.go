package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/pipeline"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultSyncFrequency = 5 * time.Minute

var (
	periodically bool
	frequency    time.Duration
	cronExpr     string
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Execute a sync pipeline",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		if cmd.Flags().Changed("source-repo") {
			repo := viper.GetString("source.repo")
			if repo != "" {
				cfg.Source.Repo = repo
			}
		}
		if cmd.Flags().Changed("dest-path") {
			path := viper.GetString("destination.path")
			if path != "" {
				cfg.Destination.Path = path
			}
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validate config: %w", err)
		}

		p, err := pipeline.Build(cfg)
		if err != nil {
			return fmt.Errorf("build pipeline: %w", err)
		}

		if cronExpr != "" {
			return runCron(cmd.Context(), p, cronExpr)
		}

		if periodically {
			return runPeriodic(cmd.Context(), p, frequency)
		}

		return p.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVar(&periodically, "periodically", false, "run sync periodically")
	syncCmd.Flags().DurationVar(&frequency, "frequency", defaultSyncFrequency, "sync frequency (e.g. 1m, 10m, 1h)")
	syncCmd.Flags().StringVar(&cronExpr, "cron", "", "cron expression for scheduling (e.g. \"*/5 * * * *\")")
}

func runPeriodic(ctx context.Context, p *pipeline.Pipeline, freq time.Duration) error {
	ticker := time.NewTicker(freq)
	defer ticker.Stop()

	slog.Info("Starting periodic sync", "frequency", freq)

	if err := p.Run(ctx); err != nil {
		slog.Error("Sync error", "err", err)
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping periodic sync")
			return nil
		case <-ticker.C:
			if err := p.Run(ctx); err != nil {
				slog.Error("Sync error", "err", err)
			}
		}
	}
}

func runCron(ctx context.Context, p *pipeline.Pipeline, expr string) error {
	c := cron.New(cron.WithParser(cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))

	entryID, err := c.AddFunc(expr, func() {
		if err := p.Run(ctx); err != nil {
			slog.Error("Sync error", "err", err)
		}
	})
	if err != nil {
		return fmt.Errorf("invalid cron expression %q: %w", expr, err)
	}

	slog.Info("Starting cron sync", "schedule", expr)
	_ = entryID

	c.Start()
	<-ctx.Done()
	slog.Info("Stopping cron sync")
	c.Stop()
	return nil
}
