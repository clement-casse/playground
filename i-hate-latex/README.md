# I Hate LaTeX ...

## ... But it is life, and sometime it is unavoidable, so that is how I limit my frustration

I sincerely do not like the LaTeX language, it is a fact that every time I use it, I end up dealing with such a massive amount of frustration that discourages me from writting !
Still, I have to admit that the quality of LaTeX documents is good and that, in research, when working on a document, it may be required to use LaTeX.
Also, in the context of my PhD, I wrote some publications in LaTeX (it was before [Typst](https://github.com/typst/typst) has been released), so being able to compile them might still be a good idea.

## The Purpose of `i-hate-latex`

This is the purpose of this experiment : Providing a Nix Flake ([the `flake.nix` file](./flake.nix)) to be able to work on a LaTeX document without having the whole latex programs polluting the actual `$PATH` of the machine.
Nix Flakes offer a pleasant experience for having a development environment and a build system that is automated and reproducible.

I decided to build my own flake to learn in the process, but I got inspired by multiple blog posts and GitHub repositories [[1][1], [2][2], [3][3], [4][4]].

I designed this flake to only be an additionnal file to existing LaTeX workspaces, so it won't spoil the experience of peoples not using Nix.
I do not provide a Nix derivation that aims to be imported by other flakes in a generic way, nor a flakes that tackles down every LaTeX documents.

## References

- [Blog Post by _Flyx_ : Exploring Nix Flakes : Build LaTeX Documents Reproducibly][1]
- [A GitHub repository also creating a Nix Flake for Building a Latex document and that could be immported in my own flake][2]
- [A GitHub repository also providing a Flake for a LaTeX document that uses the `minted` package too][3]
- [_Fasterthanlime_ article breaking down how to build and structure a Nix Flake but not for LaTeX][4]

[1]: https://flyx.org/nix-flakes-latex/
[2]: https://github.com/benide/reproducible-latex/blob/master/template/flake.nix
[3]: https://github.com/Leixb/latex-template
[4]: https://fasterthanli.me/series/building-a-rust-service-with-nix/part-10