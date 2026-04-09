package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		projects, err := st.List()
		if err != nil {
			return err
		}

		for _, p := range projects {
			fmt.Fprintln(cmd.OutOrStdout(), p)
		}

		return nil
	},
}
