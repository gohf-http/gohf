# GoHF ✨

<img align="right" width="100px" src="https://raw.githubusercontent.com/gohf-http/assets/refs/heads/main/logo.png">

[![Test](https://github.com/gohf-http/gohf/actions/workflows/test.yml/badge.svg)](https://github.com/gohf-http/gohf/actions/workflows/test.yml)
[![Go
Reference](https://pkg.go.dev/badge/github.com/gohf-http/gohf/v2.svg)](https://pkg.go.dev/github.com/gohf-http/gohf/v2)
[![Release](https://img.shields.io/github/release/gohf-http/gohf.svg?style=flat-square)](https://github.com/gohf-http/gohf/releases)
[![Documentation](https://img.shields.io/badge/documentation-yes-brightgreen.svg)](https://github.com/gohf-http/gohf#readme)
[![Maintenance](https://img.shields.io/badge/Maintained-yes-green.svg)](https://github.com/gohf-http/gohf/graphs/commit-activity)
[![License](https://img.shields.io/github/license/gohf-http/gohf)](https://github.com/gohf-http/gohf/blob/main/LICENSE)

**GO** **H**ttp **F**ramework (Golang)

# ❓ WHY GoHF

- [Easier Error Handling](#feature-easier-error-handing)
- [Middleware](#feature-middleware)
- [Sub-Router (route grouping)](#feature-sub-router)
- [Customizable Response](#feature-customizable-response)
- Lightweight
- Based on [net/http](https://pkg.go.dev/net/http)

# 📍 Getting started

```sh
go get github.com/gohf-http/gohf/v2
```

```go
import (
  "github.com/gohf-http/gohf/v2"
  "github.com/gohf-http/gohf/v2/gohf_responses"
)
```

[Hello GoHF Example](#hello-gohf-example)

# 🪄 Features

### Feature: Easier Error Handing

```go
router.Handle("GET /greeting", func(c *gohf.Context) gohf.Response {
  name := c.Req.GetQuery("name")
  if name == "" {
    return gohf_responses.NewErrorResponse(
      http.StatusBadRequest,
      errors.New("Name is required"),
    )
  }

  greeting := fmt.Sprintf("Hello, %s!", name)
  return gohf_responses.NewTextResponse(http.StatusOK, greeting)
})
```

### Feature: Middleware

```go
router.Use(func(c *gohf.Context) gohf.Response {
  token := c.Req.GetHeader("Authorization")
  if !isValidToken(token) {
    return gohf_responses.NewErrorResponse(
      http.StatusUnauthorized,
      errors.New("Invalid token"),
    )
  }

  return c.Next()
})
```

This is how middleware works in GoHF.

![middleware](https://raw.githubusercontent.com/gohf-http/assets/refs/heads/main/middleware.png)

### Feature: Sub-Router

```go
authRouter := router.SubRouter("/auth")
authRouter.Use(AuthMiddleware)
authRouter.Handle("GET /users", func(c *gohf.Context) gohf.Response {
  // ...
})
```

### Feature: Customizable Response

You can define a customizable response by implementing `gohf.Response` interface.

```go
type Response interface {
	Send(ResponseWriter, *Request)
}
```

Refer to [gohf_responses](https://github.com/gohf-http/gohf/tree/main/gohf_responses) for examples.

This is one of my favorite features, as it promotes a centralized response handler and simplifies adding additional functionality, such as logging.

# Hello GoHF Example

```go
package main

import (
  "errors"
  "fmt"
  "log"
  "net/http"

  "github.com/gohf-http/gohf/v2"
  "github.com/gohf-http/gohf/v2/gohf_responses"
)

func main() {
  router := gohf.New()

  router.Handle("GET /greeting", func(c *gohf.Context) gohf.Response {
    name := c.Req.GetQuery("name")
    if name == "" {
      return gohf_responses.NewErrorResponse(
        http.StatusBadRequest,
        errors.New("Name is required"),
      )
    }

    greeting := fmt.Sprintf("Hello, %s!", name)
    return gohf_responses.NewTextResponse(http.StatusOK, greeting)
  })

  router.Use(func(c *gohf.Context) gohf.Response {
    return gohf_responses.NewErrorResponse(
      http.StatusNotFound,
      errors.New("Page not found"),
    )
  })

  mux := router.CreateServeMux()
  log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## 🌟 Show your support

Give a ⭐️ if this project helped you!
