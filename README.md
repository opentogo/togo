# togo

This package was inspired in [gizmo from The New York Times](https://github.com/nytimes/gizmo) and it offer a simple way to build your microservices with no external dependencies in Go.

## Installation

```bash
go get github.com/opentogo/togo
```

## Examples

* They are available in the [examples](https://github.com/opentogo/togo/tree/master/examples) directory.
* There are also examples within the [GoDoc](https://godoc.org/github.com/opentogo/togo/examples).

## Log

It implements the [Apache combined log format](https://httpd.apache.org/docs/2.2/logs.html#combined) with application name as prefix.

```bash
[togo-example-service] Running at :3000
[togo-example-service] ::1 - - [10/Oct/2019:17:19:14 +0000] "GET /svc/togo HTTP/1.1" 404 19 "-" "curl/7.54.0" 0.0000
[togo-example-service] ::1 - - [10/Oct/2019:17:19:21 +0000] "GET /svc/togo/cats HTTP/1.1" 200 39 "-" "curl/7.54.0" 0.0001
```

## Contributors

- [rogeriozambon](https://github.com/rogeriozambon) Rogério Zambon - creator, maintainer
