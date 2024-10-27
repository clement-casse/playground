{
  description = "Definition of a Single-Page-Application with Go backend";

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
        # Go and Tools
        go_1_22
        golangci-lint
        gopls
        gotools
        go-tools
        delve

        # Node & Tools
        nodejs_20
        typescript
        tailwindcss
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