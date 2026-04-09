package cmd

import (
	"fmt"

	"github.com/EdgarPost/envsec/shell"
	"github.com/spf13/cobra"
)

var hookShellFlag string

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Output shell init script",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch hookShellFlag {
		case "fish":
			fmt.Fprint(cmd.OutOrStdout(), shell.FishHook)
		default:
			return fmt.Errorf("unsupported shell: %s (supported: fish)", hookShellFlag)
		}
		return nil
	},
}

func init() {
	hookCmd.Flags().StringVar(&hookShellFlag, "shell", "fish", "Shell format (fish)")
}
