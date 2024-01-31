{
  description = "Nix flake for configuring a Bevy Workspace";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    crane = {
      url = "github:ipetkov/crane";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "utils";
      };
    };
  };

  outputs = { self, nixpkgs, utils, crane, rust-overlay }: utils.lib.eachDefaultSystem (system:
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
        rustup
        rustToolchain
        pkg-config
      ];

      buildInputs = with pkgs; [
        zstd
      ] ++ lib.optionals stdenv.isLinux [
        alsa-lib
        libxkbcommon
        udev
        vulkan-loader
        wayland
        xorg.libX11
        xorg.libXcursor
        xorg.libXi
        xorg.libXrandr
      ] ++ lib.optionals stdenv.isDarwin [
        darwin.apple_sdk_11_0.frameworks.Cocoa
        rustPlatform.bindgenHook
      ];

      src = craneLib.cleanCargoSource (craneLib.path ./.);

      # Build just the cargo dependencies, so we can reuse all of that work (e.g. via cachix) when running in CI
      cargoArtifacts = craneLib.buildDepsOnly {
        inherit src;
      };

      # Common arguments for Cargo derivation have been set here to avoid repeating them later
      commonArgs = {
        inherit version src buildInputs nativeBuildInputs;
      };
    in
    with pkgs;
    {
      # formatter: Specify the formatter that will be used by the command `nix fmt`.
      formatter = nixpkgs-fmt;

      devShells.default = mkShell {
        # A Development Environment contains both dependencies for building and running the software.
        inherit buildInputs nativeBuildInputs;

        ZSTD_SYS_USE_PKG_CONFIG = true;
        LD_LIBRARY_PATH = lib.makeLibraryPath buildInputs;
      };
    });
}
