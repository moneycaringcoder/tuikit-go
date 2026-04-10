# Install

## Library

```bash
go get github.com/moneycaringcoder/tuikit-go
```

Requires Go 1.24+. No CGO required (`CGO_ENABLED=0` is safe).

## tuitest CLI

The `tuitest` binary is a thin wrapper around `go test` that adds snapshot update, JUnit/HTML reports, watch mode, and a vitest-style reporter.

=== "Homebrew"

    ```bash
    brew install moneycaringcoder/tap/tuitest
    ```

=== "Scoop"

    ```bash
    scoop bucket add moneycaringcoder https://github.com/moneycaringcoder/scoop-bucket
    scoop install tuitest
    ```

=== "Go install"

    ```bash
    go install github.com/moneycaringcoder/tuikit-go/cmd/tuitest@latest
    ```

=== "Pre-built binary"

    Download linux/darwin/windows archives (amd64 + arm64) from the
    [GitHub Releases](https://github.com/moneycaringcoder/tuikit-go/releases) page.

## Verify

```bash
tuitest --version
```

## pkg.go.dev

Full Go API reference is available at [pkg.go.dev/github.com/moneycaringcoder/tuikit-go](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go).
