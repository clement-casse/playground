{
  description = "Nix flake for creating a custom Traefik Plugin";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        # Generate a user-friendly version number for this development environement.
        #version = builtins.substring 0 8 self.lastModifiedDate;

        # Specify the version of Go for all derivétion that will use go later on.
        overlays = [
          (final: prev: {
            go = prev.go_1_22;
          })
        ];

        pkgs = import nixpkgs {
          inherit system overlays;
        };

        nativeBuildInputs = with pkgs; [
          # Go development ecosystem
          delve
          go
          gopls
          gotools
          go-tools
          golangci-lint

          yaegi
        ];
      in
      with pkgs;
      {
        # formatter: Specify the formatter that will be used by the command `nix fmt`.
        formatter = nixpkgs-fmt;

        # DevShell create a Shell with all the tools loaded to the appropriate version loaded in the $PATH
        devShells.default = mkShell {
          inherit nativeBuildInputs;
        };
      });
}
