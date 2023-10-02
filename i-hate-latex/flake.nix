{
  description = "Nix flake generating the PhD Manuscript";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    gitignore = {
      url = "github:hercules-ci/gitignore.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, utils, gitignore }:
    utils.lib.eachDefaultSystem (system:
      let
        inherit (gitignore.lib) gitignoreSource;

        pkgs = import nixpkgs { inherit system; };

        buildInputs = with pkgs; [
          coreutils
          fira-code
          fontconfig
          texlive.combined.scheme-full
          which
          python311Packages.pygments
        ];

        TEXMFHOME = ".cache";
        TEXMFVAR = ".cache/texmf-var";
        SOURCE_DATE_EPOCH = toString self.lastModified;
      in
      with pkgs;
      {
        # formatter: Specify the formatter that will be used by the command `nix fmt`.
        formatter = nixpkgs-fmt;

        devShells.default = mkShell {
          # A Development Environment contains latex, latexmk, chktex and pygmentize
          inherit buildInputs TEXMFHOME TEXMFVAR SOURCE_DATE_EPOCH;
        };

        packages.default = stdenvNoCC.mkDerivation {
          inherit buildInputs TEXMFHOME TEXMFVAR SOURCE_DATE_EPOCH;
          name = "document";
          src = gitignoreSource ./.;
          phases = [ "unpackPhase" "buildPhase" "installPhase"];

          buildPhase = ''
            runHook preBuild
            latexmk main.tex
            runHook postBuild
          '';

          installPhase = ''
            runHook preInstall
            install -m644 -D *.pdf $out/main.pdf
            runHook postInstall
          '';
        };

      });
}