package cmd

import (
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/pipeline"
	"github.com/PapaDanielVi/secret-shift/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Interactive TUI setup for secret sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := tui.Run()
		if err != nil {
			return fmt.Errorf("tui: %w", err)
		}

		p, err := pipeline.Build(cfg)
		if err != nil {
			return fmt.Errorf("build pipeline: %w", err)
		}

		return p.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
