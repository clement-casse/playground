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

        # Generate a user-friendly version number for this development environement.
        version = builtins.substring 0 8 self.lastModifiedDate;

        pkgs = import nixpkgs { inherit system; };
        nativeBuildInputs = with pkgs; [
          # Go development ecosystem
          delve
          go_1_20
          gopls
          gotools
          go-tools
          golangci-lint
          go-junit-report

          # Build Utilities
          yq-go # Inject Nix Attributes in the builder definition

          # OpenTelemetry Tools for generating a custom collector (Defined in custom derivations later)
          ocb # Generate and build the code of the custom collector
          mdatagen # Generate code for internal/metadata module in each custom modules
        ];

        # Referencing the source repository of `opentelemetry-collector` and `opentelemetry-collector-contrib`
        # to build custom tools for collector modules development.
        otelcolVersion = "0.90.0";
        otelcolSource = pkgs.fetchFromGitHub
          {
            owner = "open-telemetry";
            repo = "opentelemetry-collector";
            rev = "v${otelcolVersion}";
            sha256 = "sha256-FL0lNlofXHTdn5o6g4FHNYauzJkjCKTrspSXP9slp/A=";
          };

        otelcolContribVersion = otelcolVersion;
        otelcolContribSource = pkgs.fetchFromGitHub
          {
            owner = "open-telemetry";
            repo = "opentelemetry-collector-contrib";
            rev = "v${otelcolContribVersion}";
            sha256 = "sha256-afif2QiWoGleNoAwe5jsGDkFVj3H28uUFm01uEINBYg=";
          };

        # Define OpenTelemetry Collector Builder Binary: It does not exist in the nixpkgs repo.
        # In addition, Go binaries of OpenTelemetry Collector does not seem to be up to date.
        ocb = pkgs.buildGoModule {
          pname = "ocb"; # The Package is named `ocb` but buildGoModule installs it as `builder`
          version = otelcolVersion;
          src = otelcolSource + "/cmd/builder";
          vendorHash = "sha256-qhX5qwb/NRG8Tf2z048U6a8XysI2bJlUtUF+hfBtx4Q=";

          # Tune Build Process
          CGO_ENABLED = 0;
          ldflags = let mod = "go.opentelemetry.io/collector/cmd/builder"; in [
            "-s"
            "-w"
            "-X ${mod}/internal.version=${version}"
            "-X ${mod}/internal.date=${self.lastModifiedDate}"
          ];

          doCheck = false; # Disable running the tests on the source code (the src is external, and tests are run on the repo anyway)

          # Check that the builder is installed by asking it to display its version
          doInstallCheck = true;
          installCheckPhase = ''
            $out/bin/builder version
          '';
        };

        # Define OpenTelemetry Collector Contrib mdatagen binary: it is a binary part of the 
        # opentelemetry-collector-contrib repo to generate the `internal/metadata` package.
        mdatagen = pkgs.buildGoModule {
          pname = "mdatagen";
          version = otelcolContribVersion;
          src = otelcolContribSource + "/cmd/mdatagen";
          vendorHash = "sha256-hYTdPDsIy0MWgg3rUYNVw4LpcCnx88xyYTaw85LFHu4=";

          CGO_ENABLED = 0;
          doCheck = false;
          doInstallCheck = false;
        };
      in
      with pkgs;
      {
        # formatter: Specify the formatter that will be used by the command `nix fmt`.
        formatter = nixpkgs-fmt;

        # DevShell create a Shell with all the tools loaded to the appropriate version loaded in the $PATH
        devShells.default = mkShell {
          inherit nativeBuildInputs;
        };

        # checks: defines a collection of derivation performed to check (i.e. test, lint, ...) the code in the repository.
        checks = {
          goTest = stdenv.mkDerivation {
            pname = "UnitTests";
            inherit version nativeBuildInputs;
            src = ./.;

            configurePhase = ''
              runHook preConfigure
              export GOPATH=$NIX_BUILD_TOP/go:$GOPATH
              export GOCACHE=$TMPDIR/go-cache
              runHook postConfigure
            '';

            installPhase = ''
              runHook preInstall
              mkdir -p $out/reports
              cd ./exporter/cyphergraphexporter/
              go test -json 2>&1 ./... | go-junit-report -set-exit-code -parser gojson -out $out/reports/cyphergraph.xml
              runHook postInstall
            '';
          };
        };

        packages.default = stdenv.mkDerivation rec {
          inherit version nativeBuildInputs;
          pname = "otelcol-custom";
          src = ./.;

          outputs = [ "out" "gen" ];

          # The Patch phase modifies the source code to run with Nix:
          # In that case it retrieves the package name version and OpenTelemetry Collector Builder version
          # to inject them in the builder configuration file.
          patchPhase = ''
            runHook prePatch
            ${yq-go}/bin/yq -i '
              .dist.name = "${pname}" |
              .dist.version = "${version}" |
              .dist.otelcol_version = "${otelcolVersion}" |
              .dist.output_path = "'$gen'/go/src/${pname}"' ${builderManifestFile}
            echo "===== FILE PATCHED: ${builderManifestFile} ====="
            cat ${builderManifestFile}
            echo "================================================"
            runHook postPatch
          '';

          # The Configure phase sets the build system up for running OCB:
          # The Go environment is setup to match Nix constraints: the code generated by OCB will be send
          # in the $GO_MOD_GEN_DIR directory that is part of the GOPATH.
          configurePhase = ''
            runHook preConfigure
            mkdir -p "$gen/go/src/${pname}"
            export GOPATH=$gen/go:$GOPATH
            export GOCACHE=$TMPDIR/go-cache
            runHook postConfigure
          '';

          # Custom these values to build on specific platforms
          inherit (go_1_20) GOOS GOARCH;

          # The OCB binary is then run with the patched definition and creates the binary
          buildPhase = ''
            runHook preBuild
            ${ocb}/bin/builder --config="${builderManifestFile}"
            runHook postBuild
          '';

          # The Binary is moved from $gen to $out
          installPhase = ''
            runHook preInstall
            install -m755 -D "$gen/go/src/${pname}/${pname}" "$out/bin/${pname}"
            runHook postInstall
          '';
        };

        apps.default = { };
      });
}
