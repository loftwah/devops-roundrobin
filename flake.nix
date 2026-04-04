{
  description = "Local-first DevOps round robin training environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            bashInteractive
            curl
            docker-client
            go
            gnumake
            jq
            kind
            kubectl
            kubernetes-helm
            kustomize
            k6
            postgresql
            redis
            terraform
            terraform-docs
            tflint
            trivy
            yq-go
          ];

          shellHook = ''
            export PATH="$PWD/scripts:$PATH"
            echo "DevOps Round Robin shell ready."
            echo "Use: cp .env.example .env && make help"
          '';
        };
      }
    );
}
