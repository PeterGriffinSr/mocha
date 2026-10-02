# Mocha

## Building from source

The lexer and parser are generated from [`Mocha.g4`](internal/pkg/mochalib/parser/Mocha.g4) with [ANTLR](https://www.antlr.org/). The generated files are git-ignored, so you must generate them before your first build.

### Prerequisites

- [Go](https://go.dev/dl/) (the version in `go.mod` or newer)
- Java 11 or newer, which ANTLR runs on
- Python with `pip`, to install the ANTLR launcher

Install the launcher:

```sh
pip install antlr4-tools
```

The launcher downloads the ANTLR jar the first time it runs. Pin the version to match CI. This also avoids a failure when the launcher can't look up the latest release:

```sh
# macOS / Linux
export ANTLR4_TOOLS_ANTLR_VERSION=4.13.2

# Windows (PowerShell)
$env:ANTLR4_TOOLS_ANTLR_VERSION = "4.13.2"
```

### Build

```sh
go generate ./...                  # generate the lexer and parser
go build -o mocha ./cmd/mocha      # build the compiler (mocha.exe on Windows)
```