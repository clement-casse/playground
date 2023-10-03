{
  description = "";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        # Generate a user-friendly version number.
        version = builtins.substring 0 8 self.lastModifiedDate;

        javaVersion = 20;

        overlays = [
          (final: prev: rec {
            jdk = prev."jdk${toString javaVersion}";
            gradle = prev.gradle.override { java = jdk; };
            kotlin = prev.kotlin.override { jre = jdk; };
          })
        ];

        pkgs = import nixpkgs {
          inherit system overlays;
        };

        buildInputs = with pkgs; [
          gradle
          kotlin
        ];

      in
      with pkgs;
      {
        # formatter: Specify the formatter that will be used by the command `nix fmt`.
        formatter = nixpkgs-fmt;

        devShells.default = mkShell {
          inherit buildInputs;
        };
      });
}
