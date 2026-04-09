package resolver

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Result struct {
	ProjectKey string
	Subpath    string
}

// Resolve maps a directory to a project key and subpath.
func Resolve(dir string) (Result, error) {
	// Try git-based resolution first
	repoRoot, remote, err := gitInfo(dir)
	if err == nil && remote != "" {
		key := normalizeRemote(remote)
		subpath := ""
		if dir != repoRoot {
			rel, err := filepath.Rel(repoRoot, dir)
			if err == nil && rel != "." {
				subpath = rel
			}
		}
		return Result{ProjectKey: key, Subpath: subpath}, nil
	}

	// Fallback: path relative to ~/Code/
	home, _ := os.UserHomeDir()
	codeDir := filepath.Join(home, "Code")
	if strings.HasPrefix(dir, codeDir+"/") {
		rel, _ := filepath.Rel(codeDir, dir)
		return Result{ProjectKey: rel, Subpath: ""}, nil
	}

	return Result{}, fmt.Errorf("cannot resolve project for %s: not a git repo and not under ~/Code/", dir)
}

func gitInfo(dir string) (repoRoot string, remote string, err error) {
	// Check for worktree: use git-common-dir to find the main repo
	commonDir, err := gitCmd(dir, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", "", err
	}

	// Get the toplevel of the current working tree
	topLevel, err := gitCmd(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", "", err
	}

	// If commonDir differs from .git, this is a worktree — resolve from the main repo
	gitDir := filepath.Join(topLevel, ".git")
	commonDir = resolveCommonDir(commonDir, topLevel)

	var remoteDir string
	if commonDir != gitDir {
		// Worktree: the main repo is the parent of commonDir
		remoteDir = filepath.Dir(commonDir)
	} else {
		remoteDir = topLevel
	}

	remote, err = gitCmd(remoteDir, "remote", "get-url", "origin")
	if err != nil {
		return topLevel, "", err
	}

	return topLevel, remote, nil
}

func resolveCommonDir(commonDir, topLevel string) string {
	if filepath.IsAbs(commonDir) {
		return filepath.Clean(commonDir)
	}
	return filepath.Clean(filepath.Join(topLevel, commonDir))
}

func gitCmd(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// normalizeRemote converts git remote URLs to a canonical form.
// e.g. "git@github.com:Org/repo.git" → "github.com/Org/repo"
// e.g. "https://github.com/Org/repo.git" → "github.com/Org/repo"
func normalizeRemote(remote string) string {
	// Handle SSH URLs: git@github.com:Org/repo.git
	if strings.Contains(remote, "@") && strings.Contains(remote, ":") {
		parts := strings.SplitN(remote, "@", 2)
		if len(parts) == 2 {
			hostPath := strings.Replace(parts[1], ":", "/", 1)
			remote = "https://" + hostPath
		}
	}

	// Parse as URL
	u, err := url.Parse(remote)
	if err != nil {
		return strings.TrimSuffix(remote, ".git")
	}

	path := strings.TrimPrefix(u.Path, "/")
	path = strings.TrimSuffix(path, ".git")

	return u.Hostname() + "/" + path
}
