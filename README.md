# My Playground

[![GitHubPageBadge](https://img.shields.io/badge/GitHub%20Pages-222222?style=for-the-badge&logo=GitHub%20Pages&logoColor=white)](https://clement-casse.github.io/playground)

**This repo is my personal playground**, so please do not expect any finished or well-scoped project but rather some very minimal projects doing one thing (at most).
This is a repository where I test and publish my experiments on some technologies, languages, tools, and everything I feel the need of hacking around with.
Some of these experiments may come along with some blog-like posts, presented as a GitHup Page website hosted in this same repository.

Almost everything presented in this repository is the underlying code and implementations I will be discussing the posts.


## List of the Different Projects

In this git repository, each of the following directories correspond to one of the playground projects I have.

| Project      | Short Description           | Status |
|--------------|-----------------------------|--------|
| [`i-hate-latex`](./i-hate-latex/)   | Building a latex paper with Nix flakes, so that I can avoid installing MacTex on my machine. | Will probably not receive any update |
| [`shuttle-rust`](./shuttle-rust/)   | A Rust workspace for testing [shuttle.rs](https://www.shuttle.rs/). | Will continue later |
| [`pulumi-go`](./pulumi-go/)         | Just me discovering and soing some initial project with Pulumi in Go. | Will continue later |
| [`kotlin-jetbrains-webapp`](./kotlin-jetbrains-webapp/) | Putting my hands into Kotlin, Ktor & Exposed for the first time. | Pondering if I should continue |
| [`otelcol-custom`](./otelcol-custom/) | Creation of a custom OpenTelemetry Collector and creation of custom modules. | Active |
| [`webservice-go`](./webservice-go/) | A very conventional Web server in Go with some reusable components, I think ... | Active |

These projects come along with a `flake.nix` file that describes the development environment for this project.
The two following files are also present: 

- `.envrc`: a file used by [direnv](https://direnv.net/) to automatically load the development environment described in the Nix Flake.
- `flake.lock`: a file registering all the version used when creating the environment from the Nix Flake.

Beside these files, expect each of these directories to have a structure that fits the technology they are implemented in.

## Files for Jekyll Website

This repository also being a Jekyll website, it also have some well-known directories where Jekyll expects to find its data.

> This Jekyll site implements the [`minimal-mistakes` template](https://github.com/mmistakes/minimal-mistakes)

| Directory             | Purpose       |
|-----------------------|---------------|
| [`_data`](./_data/)   | Directory to store all Jekyll underlying data, like navigation menu and translations. |
| [`_pages`](./_pages/) | Directory where Jekyll expects to find *pages* |
| [`_posts`](./_posts/) | Directory where Jekyll expects to find the *posts* |
| [`assets`](./assets/) | Directory to store static assets used in the Jekyll Site such as images or custom scripts. |
| `_config.yml`         | Jekyll config file |
| `Gemfile`             | Gemfile for Testing the site locally |
| `Gemfile.lock`        | Lock file for tracking gem versions |

Also, the site is published on GitHub Pages, which relies on Jekyll.
All plugin versions actually deployed on GitHub Pages are references [here](https://pages.github.com/versions/).
