# A Webservice in Go

[![Go Report Card](https://goreportcard.com/badge/github.com/clement-casse/playground/webservice-go?style=flat-square)](https://goreportcard.com/report/github.com/clement-casse/playground/webservice-go)

## Why a WebService in Go into the playground ?

Well, I made some code in my work experiences where I used Go Webservices that I instrumented with OpenTelemetry.
In this project, I make an attempt to restart the design from scratch of an OpenTelemetry monitored web-based service to propose a nice architecture to be used as a base for a later application development.

## Design decisions

This projects makes an attempt to match the [Standard Go Project Layout][1].
It also replicates some code structures that I found rather elegant (_for Go code ..._) in other Open-Source Go projects:

- [Uber Go Style Guide][2] and [Golang Practical Tips][3] provide guidance to write clear and effective Go code.
- [Functional Options frome Dave Cheney][4] extensively used in [open-telemetry/opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go) where I first discovered it.
- In [Handling errors in Go web apps][5], the author proposes a clean approach to decouple the error handling from the logic of the http handlers without any additional library.
- Packages heavily tighted to a technology like `tools/users/postgres` (which implements user store on a Postgres database) are tested with [testcontainers][7].


## References

- [Standard Go Project Layout][1]
- [Uber Go Style Guide][2]
- [Golang Practical Tips][3]
- [Functional Options][4]
- [boldlygo.tech post on handling errors in Go web apps][5]
- [LogRocket article: Rate limiting your Go application][6]
- [An article on TestContainer in Go by a former colleague][7]
- [Article on implementing a password authentication in a Go application][8]

[1]: https://github.com/golang-standards/project-layout
[2]: https://github.com/uber-go/guide/blob/master/style.md
[3]: https://github.com/func25/go-practical-tips/blob/main/tips.md
[4]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
[5]: https://boldlygo.tech/posts/2024-01-08-error-handling/
[6]: https://blog.logrocket.com/rate-limiting-go-application/
[7]: https://medium.com/@nicolas.comet/go-testcontainers-4b5399b849d9
[8]: https://www.sohamkamani.com/golang/password-authentication-and-storage/
