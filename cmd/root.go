package cmd

import (
	"fmt"

	"github.com/EdgarPost/envsec/config"
	"github.com/EdgarPost/envsec/store"
	"github.com/EdgarPost/envsec/store/fs"
	"github.com/spf13/cobra"
)

var (
	cfg   config.Config
	st    store.Store
	// Set via ldflags
	version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "envsec",
	Short: "Per-directory environment variables, synced and secure",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(pathCmd)
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Fprintf(rootCmd.ErrOrStderr(), "warning: failed to load config: %v\n", err)
		cfg = config.Default()
	}
	st = fs.New(cfg.Filesystem.Path)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("envsec", version)
	},
}
