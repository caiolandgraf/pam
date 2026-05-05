{
  description = "PAM — Pam's Database Drawer, SQL query CLI tool";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "pam";
          version = "1.2.0";

          src = ./.;

          # Update by running: nix build .#default 2>&1 | grep "got:"
          vendorHash = "sha256-fjjYqn0Rui2ZfRY/eo01ZfuEK1kA3ZwO48hnqaA7V1U=";

          enableCGO = true;

          ldflags = [
            "-s"
            "-w"
            "-X main.Version=v1.2.0"
          ];

          meta = with pkgs.lib; {
            description = "Minimal CLI tool for managing SQL queries across multiple databases";
            homepage = "https://github.com/caiolandgraf/pam";
            license = licenses.mit;
            mainProgram = "pam";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            postgresql
          ];

          shellHook = ''
            echo "========================================="
            echo "PAM development environment ready!"
            echo "========================================="
            echo ""
            echo "Available tools:"
            echo "  - Go compiler"
            echo "  - PostgreSQL client (psql)"
            echo "  - SQLite client (sqlite3)"
            echo ""
          '';
        };
      }
    );
}
