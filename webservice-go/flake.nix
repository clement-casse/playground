{
  description = "Definition of a Single-Page-Application with Go backend";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
  };

  outputs = { self, nixpkgs, utils, pre-commit-hooks }: utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs {
        inherit system;
      };

      pre-commit-checks = pre-commit-hooks.lib.${system}.run {
        src = ./.;
        hooks = {
          nixpkgs-fmt.enable = true;
          golangci-lint = {
            enable = true;
            language = "golang";
          };
        };
      };

      buildInputs = with pkgs; [
        # Go and Tools
        go_1_22
        golangci-lint
        gopls
        gotools
        go-tools
        delve
        mockgen

        # Node & Tools
        nodejs_20
        typescript
        tailwindcss
      ] ++ pre-commit-checks.enabledPackages;

    in
    with pkgs;
    {
      # formatter: Specify the formatter that will be used by the command `nix fmt`.
      formatter = nixpkgs-fmt;

      checks = {
        pre-commit = pre-commit-checks;
      };

      devShells.default = mkShell {
        inherit (pre-commit-checks) shellHook;
        inherit buildInputs;
      };
    }
  );
}
