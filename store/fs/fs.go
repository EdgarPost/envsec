package fs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Store struct {
	basedir string
}

func New(basedir string) *Store {
	return &Store{basedir: basedir}
}

func (s *Store) Init(projectKey string, subpath string) error {
	path := s.envPath(projectKey, subpath)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return nil // already exists
	}
	return os.WriteFile(path, []byte{}, 0600)
}

func (s *Store) Get(projectKey string, subpath string) (map[string]string, error) {
	vars := make(map[string]string)

	// Load root-level vars first
	rootPath := s.envPath(projectKey, "")
	if err := loadEnvFile(rootPath, vars); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Load subpath-specific vars (overrides root)
	if subpath != "" {
		subPath := s.envPath(projectKey, subpath)
		if err := loadEnvFile(subPath, vars); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	return vars, nil
}

func (s *Store) Set(projectKey string, subpath string, key string, value string) error {
	path := s.envPath(projectKey, subpath)
	vars := make(map[string]string)
	if err := loadEnvFile(path, vars); err != nil && !os.IsNotExist(err) {
		return err
	}
	vars[key] = value
	return writeEnvFile(path, vars)
}

func (s *Store) Remove(projectKey string, subpath string, key string) error {
	path := s.envPath(projectKey, subpath)
	vars := make(map[string]string)
	if err := loadEnvFile(path, vars); err != nil && !os.IsNotExist(err) {
		return err
	}
	delete(vars, key)
	return writeEnvFile(path, vars)
}

func (s *Store) List() ([]string, error) {
	var projects []string
	err := filepath.Walk(s.basedir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".env") {
			rel, _ := filepath.Rel(s.basedir, path)
			// Convert file path back to project key
			// e.g. "github.com/Org/repo.env" → "github.com/Org/repo"
			project := strings.TrimSuffix(rel, ".env")
			projects = append(projects, project)
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return projects, nil
}

func (s *Store) Path(projectKey string, subpath string) ([]string, error) {
	var paths []string

	rootPath := s.envPath(projectKey, "")
	if _, err := os.Stat(rootPath); err == nil {
		paths = append(paths, rootPath)
	}

	if subpath != "" {
		subPath := s.envPath(projectKey, subpath)
		if _, err := os.Stat(subPath); err == nil {
			paths = append(paths, subPath)
		}
	}

	return paths, nil
}

func (s *Store) Import(projectKey string, subpath string, dotenvPath string) error {
	vars := make(map[string]string)
	if err := parseDotenv(dotenvPath, vars); err != nil {
		return err
	}

	path := s.envPath(projectKey, subpath)
	existing := make(map[string]string)
	if err := loadEnvFile(path, existing); err != nil && !os.IsNotExist(err) {
		return err
	}

	for k, v := range vars {
		existing[k] = v
	}

	return writeEnvFile(path, existing)
}

// envPath returns the filesystem path for a project key + subpath.
// e.g. key="github.com/Org/repo", subpath="" → "<basedir>/github.com/Org/repo.env"
// e.g. key="github.com/Org/repo", subpath="apps/backend" → "<basedir>/github.com/Org/repo/apps/backend.env"
func (s *Store) envPath(projectKey string, subpath string) string {
	if subpath == "" {
		return filepath.Join(s.basedir, projectKey+".env")
	}
	return filepath.Join(s.basedir, projectKey, subpath+".env")
}

func loadEnvFile(path string, vars map[string]string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Strip "export " prefix
		line = strings.TrimPrefix(line, "export ")

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = unquote(value)
		vars[key] = value
	}
	return scanner.Err()
}

func parseDotenv(path string, vars map[string]string) error {
	return loadEnvFile(path, vars)
}

func writeEnvFile(path string, vars map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	// Sort keys for deterministic output
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&buf, "%s=%s\n", k, quoteIfNeeded(vars[k]))
	}

	return os.WriteFile(path, []byte(buf.String()), 0600)
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func quoteIfNeeded(s string) string {
	if strings.ContainsAny(s, " \t\"'#$\\") {
		return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return s
}
