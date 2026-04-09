{
  description = "envsec — Per-directory environment variables, synced and secure";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "envsec";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-C2pSj/1UZei+Us78cCIwbz76fRtjQxFioZiIDypGgns=";
          ldflags = [ "-s" "-w" "-X github.com/EdgarPost/envsec/cmd.version=0.1.0" ];
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go gopls ];
        };
      }
    ) // {
      homeManagerModules.default = { config, lib, pkgs, ... }:
        let
          cfg = config.programs.envsec;
        in
        {
          options.programs.envsec = {
            enable = lib.mkEnableOption "envsec";

            package = lib.mkOption {
              type = lib.types.package;
              default = self.packages.${pkgs.system}.default;
              description = "The envsec package to use.";
            };

            enableFishIntegration = lib.mkEnableOption "Fish shell integration for envsec";

            storePath = lib.mkOption {
              type = lib.types.nullOr lib.types.str;
              default = null;
              example = "~/Code/envsec";
              description = "Override the storage directory. Defaults to $XDG_DATA_HOME/envsec.";
            };
          };

          config = lib.mkIf cfg.enable {
            home.packages = [ cfg.package ];

            home.sessionVariables = lib.mkIf (cfg.storePath != null) {
              ENVSEC_STORE = cfg.storePath;
            };

            programs.fish.interactiveShellInit = lib.mkIf cfg.enableFishIntegration ''
              envsec hook --shell fish | source
            '';
          };
        };
    };
}
