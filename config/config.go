package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Store      string     `toml:"store"`
	Filesystem Filesystem `toml:"filesystem"`
}

type Filesystem struct {
	Path string `toml:"path"`
}

func Default() Config {
	return Config{
		Store: "filesystem",
		Filesystem: Filesystem{
			Path: filepath.Join(homeDir(), "Code", "envsec"),
		},
	}
}

func Load() (Config, error) {
	cfg := Default()

	path := filepath.Join(configDir(), "envsec", "config.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	cfg.Filesystem.Path = expandHome(cfg.Filesystem.Path)

	return cfg, nil
}

func expandHome(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		return filepath.Join(homeDir(), path[2:])
	}
	return path
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return os.Getenv("HOME")
	}
	return home
}

func configDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir
	}
	return filepath.Join(homeDir(), ".config")
}
