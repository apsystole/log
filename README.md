# log [![Build Status](https://github.com/apsystole/log/actions/workflows/gotest.yml/badge.svg)](https://github.com/apsystole/log/actions) [![Go Reference](https://pkg.go.dev/badge/github.com/apsystole/log.svg)](https://pkg.go.dev/github.com/apsystole/log)

Go logging library for GCP App Engine, Cloud Run, Cloud Functions:

- Encodes **severity** so that coloring works in the Google web-based Logs Explorer.
- Zero **dependencies**, it is a stdlib-only module. Nul, zilch, nada!
- Safe for concurrent logging.
- Backward compatible with older Go versions as far as 1.8.

## Screenshot

![Coloring on GCP](https://i.imgur.com/KXQfr8a.png)

## Examples

```go
package main

import "github.com/apsystole/log"

func main() {
    log.Print("my message")
}
```

More details to be found in the [documentation on pkg.go.dev](https://pkg.go.dev/github.com/apsystole/log).
