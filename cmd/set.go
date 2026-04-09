package cmd

import (
	"os"

	"github.com/EdgarPost/envsec/resolver"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set KEY VALUE",
	Short: "Set an environment variable",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		result, err := resolver.Resolve(dir)
		if err != nil {
			return err
		}

		return st.Set(result.ProjectKey, result.Subpath, args[0], args[1])
	},
}
