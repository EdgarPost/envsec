package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/EdgarPost/envsec/resolver"
	"github.com/spf13/cobra"
)

var shellFlag string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Output export statements for shell integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		result, err := resolver.Resolve(dir)
		if err != nil {
			// Not in a project — output nothing, exit cleanly
			return nil
		}

		vars, err := st.Get(result.ProjectKey, result.Subpath)
		if err != nil {
			return nil // silently fail for shell hook usage
		}

		keys := make([]string, 0, len(vars))
		for k := range vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := vars[k]
			switch shellFlag {
			case "fish":
				fmt.Fprintf(cmd.OutOrStdout(), "set -gx %s %s\n", k, fishQuote(v))
			case "bash", "zsh":
				fmt.Fprintf(cmd.OutOrStdout(), "export %s=%s\n", k, shQuote(v))
			default:
				return fmt.Errorf("unsupported shell: %s", shellFlag)
			}
		}

		return nil
	},
}

func init() {
	exportCmd.Flags().StringVar(&shellFlag, "shell", "fish", "Shell format (fish, bash, zsh)")
}

func fishQuote(s string) string {
	if strings.ContainsAny(s, " \t'\"\\$#(){}|;&<>") {
		return "'" + strings.ReplaceAll(s, "'", "\\'") + "'"
	}
	return s
}

func shQuote(s string) string {
	if strings.ContainsAny(s, " \t'\"\\$#(){}|;&<>`!") {
		return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
	}
	return s
}
