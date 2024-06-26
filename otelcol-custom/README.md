# Building a Custom OpenTelemetry Collector with Nix

[![ghpages badge](https://img.shields.io/badge/GitHub%20Pages-222222?style=for-the-badge&logo=GitHub%20Pages&logoColor=white)](https://clement-casse.github.io/playground/otelcol-custom/)
[![Go Report Card](https://goreportcard.com/badge/github.com/clement-casse/playground/otelcol-custom)](https://goreportcard.com/report/github.com/clement-casse/playground/otelcol-custom)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/clement-casse/playground/otlecol-custom.yml "Build")](https://github.com/clement-casse/playground/actions/workflows/otlecol-custom.yml)

## Why `otelcol-custom` ?

This project aims to demonstrate how to create custom OpenTelemetry Collectors with custom modules by leveraging the tools and methods used by the OpenTelemetry project.
The structure of this project mimics the struture of [the OpenTelemetry Collector Contrib repo](https://github.com/open-telemetry/opentelemetry-collector-contrib) to arrange the module in their own folders.
Still, with this project I intend to provide an example of an entry level contribution on OpenTelemetry-collector module that would not aim to be merged with the core project.

So, in order to lower complexity to a more familiar level for occasionnal Go developper like me, I took some distance with some of the practices used in the OpenTelemetry that I personaly am not so big of a fan of, like:

- I do not run the scripts through crypic [Makefiles](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/Makefile): I tend to minimize the reliability on Makefiles to only use Go base tooling, I prefer using commands like `go generate ./...` and that run script from Go magic comments and if special tools are required (like [`mdatagen`](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/mdatagen)) I invoke them from `go run` commands, [reference here](https://www.jvt.me/posts/2022/06/15/go-tools-dependency-management/).
- I do not use [chloggen](https://go.opentelemetry.io/build-tools/chloggen) no any of the [opentelemetry-go-build-tools](https://github.com/open-telemetry/opentelemetry-go-build-tools) (yet ?): my project is not at a maturity level where these tools would be benefical.
- The project is built with Nix, althought it is not a hard dependancy and can be omited (and should by not conoiseurs), it allows to unify the build process and the required tooling. While Nix may have high potential on build reproductibility, I am not familiar enougth with the tool to use it in CI too, but I am looking forward to it, currently it is only used to provide a dev shell with Go and the ocb binary and to run ocb with on the [./builder-config.yaml](./builder-config.yaml) file.
- a Dockerfile defines the image specification of this custom OpenTelemetry container that can run all sub directories in the [`./examples/`](./examples/) directory.
- CI/CD is lighter and more leniant, it is defined as part of the [playground repository](../.github/workflows/otlecol-custom.yml).

## Content

| Module name | Type of module | Status      | Description           |
|:------------|:--------------:|:-----------:|:----------------------|
| cyphergraph | exporter       | In progress | An exporter that send traces to a Neo4j compatible database and encoding traces in a hierarchical property graph |

## Project Status

- [x] Build a Custom Collector with Nix Flakes (#2)
- [x] Integrate a Skeleton of exporter module to add to this specific Collector build (#4)
- [ ] Add the exporter logic

## References

1. [Some blog explaining Nix to write derivations][1]
2. [Internal of Nix to create an environment to build Go Modules][2]
3. [A Collection of Articles about learning Nix][3]

[1]: https://blog.ysndr.de/posts/internals/2021-01-01-flake-ification/
[2]: https://github.com/NixOS/nixpkgs/blob/e3fbbb1d108988069383a78f424463e6be087707/pkgs/development/go-packages/generic/default.nix#L92-L110
[3]: https://ianthehenry.com/posts/how-to-learn-nix/
