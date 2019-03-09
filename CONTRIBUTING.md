# Contributing

Hello there! Thanks for willing to help me improve Magnanimous!

All ideas are welcome. But to keep the project simple and focused on its main goals, not all contributions may be
accepted.

Please create a issue on GitHub first before you put too much work on an idea that may not be accepted.

## File structure

The root folder contains the `main.go` file which is the Magnanimous CLI runner, but most source code is put inside
the `mg` directory (and Go package) and sub-directories.

The `tests` directory contains Go tests.

The `website` directory contains the [Magnanimous Website](https://renatoathaydes.github.io/magnanimous) source code.

The `sample` directory contains sample projects using Magnanimous.

## Compiling and testing

Magnanimous is written in the Go Programming Language.

It uses the new [Go module](https://github.com/golang/go/wiki/Modules) system (see [go.mod](go.mod))
to manage dependencies.

To compile it:

```bash
go build
```

To test:

```bash
go test ./...
```

To upgrade dependencies:

```bash
go get -u
```

### Tools

I personally use IntelliJ Ultimate for development, but the project is just standard Go and any Go IDE would do.

The website is written using pure HTML, CSS (no JS) and, of course, Magnanimous!

Website graphics are being created online with [Figma](https://www.figma.com/file/WWNwFQocI5vDQd2pdwJzCHJ3/magnanimous-transformation?node-id=0%3A1).

This project is hosted on [GitHub](https://github.com/).

## Creating pull requests

To create a pull request, use [GitHub](https://help.github.com/en/articles/about-pull-requests).

Please target the `next` branch, not `master`. If the branch doesn't exist (I keep deleting it after merging to master)
please create an issue and I'll add it back.