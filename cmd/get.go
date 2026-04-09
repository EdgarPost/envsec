package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/EdgarPost/envsec/resolver"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [KEY]",
	Short: "Show environment variables",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		result, err := resolver.Resolve(dir)
		if err != nil {
			return err
		}

		vars, err := st.Get(result.ProjectKey, result.Subpath)
		if err != nil {
			return err
		}

		if len(args) == 1 {
			val, ok := vars[args[0]]
			if !ok {
				return fmt.Errorf("variable %s not found", args[0])
			}
			fmt.Fprintln(cmd.OutOrStdout(), val)
			return nil
		}

		keys := make([]string, 0, len(vars))
		for k := range vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, vars[k])
		}

		return nil
	},
}
