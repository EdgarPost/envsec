package store

type Store interface {
	// Get returns all env vars for a project key + subpath (with inheritance).
	Get(projectKey string, subpath string) (map[string]string, error)
	// Set sets a single env var at the given project key + subpath level.
	Set(projectKey string, subpath string, key string, value string) error
	// Remove removes a single env var.
	Remove(projectKey string, subpath string, key string) error
	// List returns all registered project keys.
	List() ([]string, error)
	// Init creates the env file for a project key + subpath.
	Init(projectKey string, subpath string) error
	// Path returns the resolved file path(s) for the project.
	Path(projectKey string, subpath string) ([]string, error)
	// Import reads a dotenv file and writes all vars into the store.
	Import(projectKey string, subpath string, dotenvPath string) error
}
