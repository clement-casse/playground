{
  description = "Nix flake for creating a custom OpenTelemetry Collector";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        # Relative path of the config file describing the modules embedded in the custom OpenTelemetry Collector.
        builderManifestFile = "builder-config.yaml";

        # Generate a user-friendly version number.
        version = builtins.substring 0 8 self.lastModifiedDate;

        pkgs = import nixpkgs { inherit system; };

        # Define OpenTelemetry Collector Builder Binary build It does not exist in the nixpkgs repo.
        # In addition, Go binaries of OpenTelemetry Collector does not seem to be up to date.
        ocb = pkgs.buildGoModule rec {
          pname = "ocb"; # The Package is named `ocb` but buildGoModule installs it as `builder`
          version = "0.86.0";

          src = pkgs.fetchFromGitHub
            {
              owner = "open-telemetry";
              repo = "opentelemetry-collector";
              rev = "v${version}";
              sha256 = "sha256-Ucp00OjyPtHA6so/NOzTLtPSuhXwz6A2708w2WIZb/E=";
            } + "/cmd/builder";
          vendorHash = "sha256-MTwD9xkrq3EudppLSoONgcPCBWlbSmaODLH9NtYgVOk=";

          # Tune Build Process: Set Go LDFlags and Disable CGo
          # TODO: Address the log "You're building a distribution with non-aligned version of the builder.
          #       Compilation may fail due to API changes."
          CGO_ENABLED = 0;
          ldflags = let mod = "go.opentelemetry.io/collector/cmd/builder"; in [
            "-s"
            "-w"
            "-X ${mod}/internal.version=${version}"
            "-X ${mod}/internal.date=${self.lastModifiedDate}"
          ];

          # The go test command fails in Nix Env with the following Error:
          # > failed to update go.mod: exit status 1. Output:
          # > go.opentelemetry.io/collector/cmd/builder imports
          # >      go.opentelemetry.io/collector/component: github.com/knadh/koanf/maps@v0.1.1: module lookup disabled by GOPROXY=off
          # TODO: Address this issue with less aggressive test disabling or with system configuration ?
          doCheck = false;

          # Check that the builder is installed by asking it to display its version
          doInstallCheck = true;
          installCheckPhase = ''
            $out/bin/builder version
          '';
        };

        buildInputs = with pkgs; [
          go_1_20
          gopls
          gotools
          go-tools
          ocb
          yq-go
        ];

      in
      with pkgs;
      {
        # formatter: Specify the formatter that will be used by the command `nix fmt`.
        formatter = nixpkgs-fmt;

        devShells.default = mkShell {
          inherit buildInputs;
        };

        packages.default = stdenv.mkDerivation rec {
          inherit version buildInputs;
          pname = "otelcol-custom";

          src = self;

          inherit (go) GOOS GOARCH;

          patchPhase = ''
            runHook prePatch
            ${yq-go}/bin/yq -i '
              .dist.name = "${pname}" |
              .dist.version = "${version}" |
              .dist.otelcol_version = "${ocb.version}" |
              .dist.output_path = "'$NIX_BUILD_TOP/go/src/$(dirname "${self}")'"' ${builderManifestFile}
            echo "===== FILE PATCHED: ${builderManifestFile} ====="
            cat ${builderManifestFile}
            echo "================================================"
            runHook postPatch
          '';

          configurePhase = ''
            runHook preConfigure
            mkdir -p "$NIX_BUILD_TOP/go/src/$(dirname "${self}")"
            export GOPATH=$NIX_BUILD_TOP/go:$GOPATH
            export GOCACHE=$TMPDIR/go-cache
            runHook postConfigure
          '';

          buildPhase = ''
            runHook preBuild
            ${ocb}/bin/builder --config="${builderManifestFile}"
            runHook postBuild
          '';

          installPhase = ''
            runHook preInstall
            cd "$NIX_BUILD_TOP/go/src/$(dirname "${self}")"
            mkdir -p $out/bin
            install -m755 -D ${pname} $out/bin/${pname}
            runHook postInstall
          '';
        };

        apps.default = { };
      });
}
