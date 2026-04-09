# envsec

Per-directory environment variables, synced and secure.

envsec maps directories to environment variables using git remotes as project identifiers. Variables are stored as plain `.env` files, making them easy to sync across machines with tools like Syncthing.

## Install

### Nix flake (recommended)

Add to your flake inputs:

```nix
inputs.envsec.url = "github:EdgarPost/envsec";
```

Enable via home-manager:

```nix
imports = [ inputs.envsec.homeManagerModules.default ];

programs.envsec = {
  enable = true;
  enableFishIntegration = true;
  # storePath = "~/Code/envsec";  # optional, defaults to ~/.local/share/envsec
};
```

### From source

```bash
go install github.com/EdgarPost/envsec@latest
```

## Usage

```bash
cd ~/Code/github.com/Org/my-project

envsec init                          # register project (auto-detects git remote)
envsec set DATABASE_URL "postgres://localhost/mydb"
envsec set API_KEY "sk-123"
envsec get                           # list all vars
envsec get DATABASE_URL              # get a single var
envsec rm API_KEY                    # remove a var
envsec import .env.local             # import from an existing .env file
envsec edit                          # open env file in $EDITOR
envsec path                          # print env file path(s)
envsec list                          # list all registered projects
```

With the Fish hook active, environment variables are automatically loaded when you `cd` into a project and unloaded when you leave. Any mutation (`set`, `rm`, `import`) reloads vars immediately.

## Shell integration

### Fish

```fish
envsec hook --shell fish | source
```

Or add to your Fish config / home-manager (the Nix module does this for you with `enableFishIntegration`).

## How it works

### Project resolution

envsec resolves the current directory to a project key using:

1. **Git remote** — `git remote get-url origin` normalized to `github.com/Org/repo`
2. **Worktrees** — uses `git rev-parse --git-common-dir` so all worktrees of the same repo share vars
3. **Fallback** — path relative to `~/Code/`

### Storage

Env files are stored in the data directory (default `~/.local/share/envsec/`):

```
~/.local/share/envsec/
  github.com/
    Org/
      repo.env                        # root-level vars
      repo/
        apps/
          backend.env                 # subpath-specific vars
          frontend.env
```

### Inheritance

In a monorepo, subpath vars inherit from and override root vars:

```bash
cd ~/Code/github.com/Org/repo
envsec set DATABASE_URL "postgres://shared"

cd apps/backend
envsec init
envsec set API_KEY "sk-backend"
envsec get
# DATABASE_URL=postgres://shared    (inherited from root)
# API_KEY=sk-backend                (subpath-specific)
```

## Configuration

### Storage path

The storage directory can be configured in three ways (highest priority first):

1. **Environment variable**: `ENVSEC_STORE=~/my/path`
2. **Config file**: `~/.config/envsec/config.toml`
3. **Default**: `$XDG_DATA_HOME/envsec` (usually `~/.local/share/envsec`)

Config file example:

```toml
[filesystem]
path = "~/Code/envsec"
```

With the Nix home-manager module:

```nix
programs.envsec = {
  enable = true;
  storePath = "~/Code/envsec";
};
```

## Syncing

The env files are plain `KEY=value` text files. Sync the storage directory between machines using any file sync tool (Syncthing, rsync, Dropbox, etc.). Files are created with `0600` permissions.

## License

MIT
