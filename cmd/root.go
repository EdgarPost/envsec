package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/EdgarPost/envsec/store"
	"github.com/EdgarPost/envsec/store/fs"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var (
	st store.Store
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
	cobra.OnInitialize(initStore)

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

type config struct {
	Store      string `toml:"store"`
	Filesystem struct {
		Path string `toml:"path"`
	} `toml:"filesystem"`
}

func loadConfig() config {
	home, _ := os.UserHomeDir()
	cfg := config{Store: "filesystem"}
	cfg.Filesystem.Path = filepath.Join(home, "Code", "envsec")

	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = filepath.Join(home, ".config")
	}

	data, err := os.ReadFile(filepath.Join(configDir, "envsec", "config.toml"))
	if err != nil {
		return cfg
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg
	}

	// Expand ~ in path
	if len(cfg.Filesystem.Path) >= 2 && cfg.Filesystem.Path[:2] == "~/" {
		cfg.Filesystem.Path = filepath.Join(home, cfg.Filesystem.Path[2:])
	}

	return cfg
}

func initStore() {
	cfg := loadConfig()
	st = fs.New(cfg.Filesystem.Path)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("envsec", version)
	},
}
