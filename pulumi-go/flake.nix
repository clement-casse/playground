{
  description = "Nix flake for testing Pulumi with Go";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "utils";
      };
    };
  };

  outputs = { self, nixpkgs, utils, gomod2nix }:
    utils.lib.eachDefaultSystem (system:
      let
        # Generate a user-friendly version number.
        version = builtins.substring 0 8 self.lastModifiedDate;

        overlays = [
          gomod2nix.overlays.default
        ];

        pkgs = import nixpkgs {
          inherit system overlays;
        };

        buildInputs = with pkgs; [
          go_1_21
          gopls
          gotools
          go-tools
          gomod2nix.packages.${system}.default
          pulumi-bin
        ];

      in
      with pkgs;
      {
        # formatter: Specify the formatter that will be used by the command `nix fmt`.
        formatter = nixpkgs-fmt;

        packages.default = buildGoModule {
          pname = "prepare-k8s-cluster";
          inherit version;
          src = ./.;
          subPackages = [ "prepare-k8s-cluster" ];

          vendorSha256 = "sha256-AMdPADIN5bgL5OZBu0SqI/MEUiN8mUZm1QJh+l6TzYA=";
          #modules = ./gomod2nix.toml;
        };

        devShells.default = mkShell {
          inherit buildInputs;
        };
      });
}
