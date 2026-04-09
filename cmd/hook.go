package cmd

import (
	"fmt"
	"os"

	"github.com/EdgarPost/envsec/shell"
	"github.com/spf13/cobra"
)

var hookShellFlag string

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Output shell init script",
	RunE: func(cmd *cobra.Command, args []string) error {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("could not determine envsec path: %w", err)
		}

		switch hookShellFlag {
		case "fish":
			fmt.Fprint(cmd.OutOrStdout(), shell.FishHook(exe))
		default:
			return fmt.Errorf("unsupported shell: %s (supported: fish)", hookShellFlag)
		}
		return nil
	},
}

func init() {
	hookCmd.Flags().StringVar(&hookShellFlag, "shell", "fish", "Shell format (fish)")
}
