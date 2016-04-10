## Package: environment

[![](https://godoc.org/github.com/yogin/go-packages/environment?status.svg)](http://godoc.org/github.com/yogin/go-packages/environment)

This package allows the use of JSON configuration files for different environments.

Usage:
```go
package main

import (
  "github.com/yogin/go-packages/environment"
  "fmt"
)

func main() {
  // Use GO_ENV variable to set the environment
  environment.Init()
  // Or you can pass it in:
  // environment.Init("production")

  // Then you can access the environment and it's configuration
  fmt.Printf("Environment %s : %+v\n", environment.Name(), environment.Get())
}
```

Environments are whitelisted to prevent unwanted injections, and includes:
 
 * `test` 
 * `development`
 * `staging`
 * `production`

You can however register new environments for your own needs:

```go
func main() {
  environment.Register("qa")
  environment.Init("qa")
  ...
}
```

Each environment has a JSON configuration file associated with it:

environment | path
---|---
`test` |  `environments/test.json`
`development` | `environments/development.json`
`staging` | `environments/staging.json`
`production` | `environments/production.json`
`...` | `...`

