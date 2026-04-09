package cmd

import (
	"fmt"
	"os"

	"github.com/EdgarPost/envsec/resolver"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import FILE",
	Short: "Import variables from a dotenv file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		result, err := resolver.Resolve(dir)
		if err != nil {
			return err
		}

		if err := st.Import(result.ProjectKey, result.Subpath, args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Imported %s\n", args[0])
		return nil
	},
}
