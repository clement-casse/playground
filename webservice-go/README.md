# A Webservice in Go

[![Go Report Card](https://goreportcard.com/badge/github.com/clement-casse/playground/webservice-go?style=flat-square)](https://goreportcard.com/report/github.com/clement-casse/playground/webservice-go)

## Why a WebService in Go into the playground ?

Well, I made some code in my work experiences where I used Go Webservices that I instrumented with OpenTelemetry.
In this project, I make an attempt to restart the design from scratch of an OpenTelemetry monitored web-based service to propose a nice architecture to be used as a base for a later application development.

## Design decisions

This projects makes an attempt to match the [Standard Go Project Layout][1].
It also replicates some code structures that I found rather elegant (_for Go code ..._) in other Open-Source Go projects:

- [Functional Options frome Dave Cheney][2] extensively used in [open-telemetry/opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go) where I first discovered it.
- In [Handling errors in Go web apps][3], the author proposes a clean approach to decouple the error handling from the logic of the http handlers without any additional library.


## References

- [Standard Go Project Layout][1]
- [Functional Options][2]
- [boldlygo.tech post on handling errors in Go web apps][3]
- [LogRocket article: Rate limiting your Go application][4]

[1]: https://github.com/golang-standards/project-layout
[2]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
[3]: https://boldlygo.tech/posts/2024-01-08-error-handling/
[4]: https://blog.logrocket.com/rate-limiting-go-application/
