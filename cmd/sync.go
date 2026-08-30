package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/pipeline"
	"github.com/PapaDanielVi/secret-shift/internal/server"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultSyncFrequency = 5 * time.Minute

var (
	periodically bool
	frequency    time.Duration
	cronExpr     string
	dryRun       bool
	serverMode   bool
	healthPort   int
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Execute a sync pipeline",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(viper.GetViper(), cfgFile)
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

		cfg.DryRun = dryRun

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validate config: %w", err)
		}

		p, err := pipeline.Build(cmd.Context(), cfg)
		if err != nil {
			return fmt.Errorf("build pipeline: %w", err)
		}

		if serverMode {
			return runServer(cmd.Context(), p, cfg, healthPort)
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
	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "simulate sync without writing to destination")
	syncCmd.Flags().BoolVar(&serverMode, "server", false, "start HTTP server with health endpoints")
	syncCmd.Flags().IntVar(&healthPort, "health-port", 8080, "port for health HTTP server")
}

func runPeriodic(ctx context.Context, p *pipeline.Pipeline, freq time.Duration) error {
	ticker := time.NewTicker(freq)
	defer ticker.Stop()

	slog.Info("Starting periodic sync", "frequency", freq)

	if err := p.Run(ctx); err != nil {
		slog.Error("sync error", "err", err)
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping periodic sync")
			return nil
		case <-ticker.C:
			if err := p.Run(ctx); err != nil {
				slog.Error("sync error", "err", err)
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
			slog.Error("sync error", "err", err)
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

func runServer(ctx context.Context, p *pipeline.Pipeline, _ *config.Config, port int) error {
	health := server.NewHealthServer(port)

	go runServerSync(ctx, p.Run, health)

	return health.Start(ctx)
}

func runServerSync(ctx context.Context, run func(context.Context) error, health *server.HealthServer) {
	report := func() {
		if err := run(ctx); err != nil {
			health.ReportSyncError(err)
		} else {
			health.ReportSyncSuccess()
		}
	}

	report()

	ticker := time.NewTicker(defaultSyncFrequency)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			report()
		}
	}
}
