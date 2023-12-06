# Building a Custom OpenTelemetry Collector with Nix

[![ghpages badge](https://img.shields.io/badge/GitHub%20Pages-222222?style=for-the-badge&logo=GitHub%20Pages&logoColor=white)](https://clement-casse.github.io/playground/otelcol-custom/)

## Why `otelcol-custom` ?

Because I work on OpenTelemetry Collector and I made a similar approach at work, without Nix Flakes and with CUE.
I have never been happy with that build process, and, therefore I made this experiment of building a collector with a Nix Flake.

The resulting flake was quite cumbersome to build and I dived into some heavy Nix Blog posts [[1], [3]] and I even had to looked at the Nix Source code [[2]] to understand how Go Modules were built (it still mysterious though).

## Project Status

- [x] Build a Custom Collector with Nix Flakes (#2)
- [x] Integrate a Skeleton of exporter module to add to this specific Collector build (#4)
- [ ] Add the exporter logic
- [ ] Run the tests in Nix flakes `checks`
- [ ] ...

## References

- [Some blog explaining Nix to write derivations][1]
- [Internal of Nix to create an environment to build Go Modules][2]
- [A Collection of Articles about learning Nix][3]

[1]: https://blog.ysndr.de/posts/internals/2021-01-01-flake-ification/
[2]: https://github.com/NixOS/nixpkgs/blob/e3fbbb1d108988069383a78f424463e6be087707/pkgs/development/go-packages/generic/default.nix#L92-L110
[3]: https://ianthehenry.com/posts/how-to-learn-nix/
