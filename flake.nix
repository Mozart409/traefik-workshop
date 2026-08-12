{
  description = "Traefik Workshop - Development Environment";

  inputs = {
    nixpkgs-linux.url = "github:NixOS/nixpkgs/nixos-unstable";
    nixpkgs-darwin.url = "github:NixOS/nixpkgs/nixpkgs-26.05-darwin";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs-linux,
    nixpkgs-darwin,
    flake-utils,
  }: let
    pkgsFor = system:
      if system == "x86_64-linux" || system == "aarch64-linux"
      then nixpkgs-linux.legacyPackages.${system}
      else if system == "x86_64-darwin" || system == "aarch64-darwin"
      then nixpkgs-darwin.legacyPackages.${system}
      else nixpkgs-linux.legacyPackages.${system};
  in
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = pkgsFor system;
        isDarwin = system == "x86_64-darwin" || system == "aarch64-darwin";
      in {
        devShells.default = pkgs.mkShell {
          name = "traefik-workshop";

          packages = with pkgs;
            [
              # Container tools - core requirement for this workshop
              docker
              docker-compose
            ]
            ++ (
              if !isDarwin
              then [podman podman-compose]
              else []
            )
            ++ [
              cocogitto
              keep-sorted
              lefthook
              # Kubernetes tools
              kubectl
              kind
              minikube
              kubernetes-helm

              # Go - for the fiboService application
              go

              # Node.js - for the Next.js application
              nodejs_22
              pnpm

              # Traefik - the main focus of this workshop
              traefik

              # Build tools mentioned in README
              gnumake

              # Code formatting (project has dprint.json)
              alejandra
              dprint
              yamllint

              # API testing and debugging utilities
              curl
              jq
              httpie
              wget

              # Git (always useful)
              git

              # Additional utilities
              watch
              tree
              bat
              fd
              ripgrep
            ];

          shellHook = ''
            echo "=================================================="
            echo "  Traefik Workshop Development Environment"
            echo "=================================================="
            echo ""
            echo "Available tools:"
            echo "  - Docker & Docker Compose"
            echo "  - Podman & Podman Compose (Docker alternative)"
            echo "  - Kubernetes: kubectl, kind, minikube, helm"
            echo "  - Go (for fiboService)"
            echo "  - Node.js 22 (for Next.js app)"
            echo "  - Traefik"
            echo "  - Make"
            echo "  - dprint (code formatter)"
            echo ""
            echo "Workshop folders:"
            echo "  01-native-clis/      - Native CLI comparison"
            echo "  02-docker-run/       - Docker run examples"
            echo "  03-docker-compose/   - Docker Compose examples"
            echo "  04-traefik/          - Basic Traefik setup"
            echo "  05-traefik-advanced/ - Advanced Traefik config"
            echo ""
            echo "=================================================="

            lefthook install 2>/dev/null || true
          '';
        };
      }
    );
}
