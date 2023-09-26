{
  # This Nix flake is heavily inspired (i.e. copied) by fasterthanlime's flake described in
  # https://fasterthanli.me/series/building-a-rust-service-with-nix/part-10
  description = "Nix flake for testing Shuttle.rs";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "utils";
      };
    };
    crane = {
      url = "github:ipetkov/crane";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "utils";
        rust-overlay.follows = "rust-overlay";
      };
    };
  };

  outputs = { self, nixpkgs, utils, rust-overlay, crane }:
    # `eachDefaultSystem` generates an output that simply has `packages.default` instead of packages."<system>".default
    utils.lib.eachDefaultSystem (system:
      let
        # Generate a user-friendly version number.
        version = builtins.substring 0 8 self.lastModifiedDate;

        overlays = [
          (import rust-overlay)
        ];

        pkgs = import nixpkgs {
          inherit system overlays;
        };

        # Get Rust Toolchain version from ./rust-toolchain.toml
        rustToolchain = pkgs.pkgsBuildHost.rust-bin.fromRustupToolchainFile ./rust-toolchain.toml;
        craneLib = (crane.mkLib pkgs).overrideToolchain rustToolchain; # make crane aware of the custom rust toolchain: https://github.com/ipetkov/crane/blob/8b08e96c9af8c6e3a2b69af5a7fa168750fcf88e/examples/cross-musl/flake.nix#L35

        # Dependencies needed at compile-time
        nativeBuildInputs = with pkgs; [
          rustToolchain
          pkg-config
        ];

        # Dependencies needed at run-time
        buildInputs = with pkgs; [
          cargo-shuttle
          docker
          openssl
        ] ++ lib.optionals stdenv.isDarwin [
          darwin.apple_sdk.frameworks.Security
          darwin.apple_sdk.frameworks.IOKit
        ];

        # Common arguments for Cargo derivation have been set here to avoid repeating them later
        commonArgs = {
          inherit version buildInputs nativeBuildInputs;
          src = craneLib.cleanCargoSource (craneLib.path ./.);
        };

        # Build just the cargo dependencies, so we can reuse all of that work (e.g. via cachix) when running in CI
        cargoArtifacts = craneLib.buildDepsOnly (commonArgs // {
          pname = "shuttle-playground-workspace";
        });
      in
      with pkgs;
      {
        # formatter: Specify the formatter that will be used by the command `nix fmt`.
        formatter = nixpkgs-fmt;

        # checks: define tests and linters to be run uppon `nix flake check`.
        checks = rec {
          default = clippy;

          clippy = craneLib.cargoClippy {
            inherit cargoArtifacts;
            cargoClippyExtraArgs = "--all-targets -- --deny warnings";
          };

          cargoFmt = craneLib.cargoFmt cargoArtifacts;
        };

        # devShell: Define Development Shell that will load all dependencies for developping and building the system.
        devShells.default = mkShell {
          # A Development Environment contains both dependencies for building and running the software.
          inherit buildInputs nativeBuildInputs;
        };

      }
    );
}
