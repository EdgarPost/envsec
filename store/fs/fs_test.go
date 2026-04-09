package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesEnvFile(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	if err := s.Init("github.com/Org/repo", ""); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "github.com", "Org", "repo.env")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected env file at %s: %v", path, err)
	}
}

func TestInitSubpath(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	if err := s.Init("github.com/Org/repo", "apps/backend"); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "github.com", "Org", "repo", "apps", "backend.env")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected env file at %s: %v", path, err)
	}
}

func TestSetAndGet(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	s.Init("github.com/Org/repo", "")
	if err := s.Set("github.com/Org/repo", "", "FOO", "bar"); err != nil {
		t.Fatal(err)
	}

	vars, err := s.Get("github.com/Org/repo", "")
	if err != nil {
		t.Fatal(err)
	}

	if vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got FOO=%s", vars["FOO"])
	}
}

func TestInheritance(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	// Set root-level var
	s.Init("github.com/Org/repo", "")
	s.Set("github.com/Org/repo", "", "DATABASE_URL", "postgres://root")
	s.Set("github.com/Org/repo", "", "SHARED", "from-root")

	// Set subpath var (overrides SHARED)
	s.Init("github.com/Org/repo", "apps/backend")
	s.Set("github.com/Org/repo", "apps/backend", "SHARED", "from-sub")
	s.Set("github.com/Org/repo", "apps/backend", "API_KEY", "secret")

	vars, err := s.Get("github.com/Org/repo", "apps/backend")
	if err != nil {
		t.Fatal(err)
	}

	if vars["DATABASE_URL"] != "postgres://root" {
		t.Errorf("expected inherited DATABASE_URL, got %s", vars["DATABASE_URL"])
	}
	if vars["SHARED"] != "from-sub" {
		t.Errorf("expected SHARED=from-sub (override), got %s", vars["SHARED"])
	}
	if vars["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %s", vars["API_KEY"])
	}
}

func TestRemove(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	s.Init("github.com/Org/repo", "")
	s.Set("github.com/Org/repo", "", "FOO", "bar")
	s.Set("github.com/Org/repo", "", "BAZ", "qux")

	if err := s.Remove("github.com/Org/repo", "", "FOO"); err != nil {
		t.Fatal(err)
	}

	vars, _ := s.Get("github.com/Org/repo", "")
	if _, ok := vars["FOO"]; ok {
		t.Error("expected FOO to be removed")
	}
	if vars["BAZ"] != "qux" {
		t.Error("expected BAZ to still be present")
	}
}

func TestImport(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)
	s.Init("github.com/Org/repo", "")

	// Create a dotenv file to import
	dotenv := filepath.Join(dir, ".env.local")
	content := `# Comment
FOO=bar
BAZ="quoted value"
SINGLE='single quoted'
export EXPORTED=yes

SPACES=  spaced
`
	os.WriteFile(dotenv, []byte(content), 0600)

	if err := s.Import("github.com/Org/repo", "", dotenv); err != nil {
		t.Fatal(err)
	}

	vars, _ := s.Get("github.com/Org/repo", "")

	tests := map[string]string{
		"FOO":      "bar",
		"BAZ":      "quoted value",
		"SINGLE":   "single quoted",
		"EXPORTED": "yes",
		"SPACES":   "spaced",
	}

	for k, want := range tests {
		if got := vars[k]; got != want {
			t.Errorf("%s: expected %q, got %q", k, want, got)
		}
	}
}

func TestList(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	s.Init("github.com/Org/repo", "")
	s.Init("github.com/EdgarPost/project", "")

	projects, err := s.List()
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
}

func TestPath(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	s.Init("github.com/Org/repo", "")
	s.Init("github.com/Org/repo", "apps/backend")

	paths, err := s.Path("github.com/Org/repo", "apps/backend")
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d: %v", len(paths), paths)
	}
}
