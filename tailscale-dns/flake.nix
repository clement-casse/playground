{
  description = "Definition of a dev environment for Fly.io deployment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }: utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs {
        inherit system;
      };

      buildInputs = with pkgs; [
        # fly.io tool
        flyctl

        # pulumi tool
        pulumi-bin

        # Go and Tools
        go_1_22
        golangci-lint
        gopls
        gotools
        go-tools
        delve
      ];
    in
    with pkgs;
    {
      # formatter: Specify the formatter that will be used by the command `nix fmt`.
      formatter = nixpkgs-fmt;

      devShells.default = mkShell {
        inherit buildInputs;
      };
    }
  );
}