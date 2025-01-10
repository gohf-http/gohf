# GoHF ✨

<img align="right" width="100px" src="https://raw.githubusercontent.com/gohf-http/assets/refs/heads/main/logo.png">

[![build](https://github.com/gohf-http/gohf/actions/workflows/test.yml/badge.svg)](https://github.com/gohf-http/gohf/actions/workflows/test.yml)
[![Go
Reference](https://pkg.go.dev/badge/github.com/gohf-http/gohf.svg)](https://pkg.go.dev/github.com/gohf-http/gohf)
[![Release](https://img.shields.io/github/release/gohf-http/gohf.svg?style=flat-square)](https://github.com/gohf-http/gohf/releases)
[![Documentation](https://img.shields.io/badge/documentation-yes-brightgreen.svg)](https://github.com/gohf-http/gohf#readme)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/gohf-http/gohf/graphs/commit-activity)
[![License](https://img.shields.io/github/license/gohf-http/gohf)](https://github.com/gohf-http/gohf/blob/main/LICENSE)

**GO** **H**ttp **F**ramework

# ❓ WHY GoHF

- [Easier error handling](#feature-easier-error-handing)
- [Support middleware](#feature-middleware)
- [Support sub-routers (route grouping)](#feature-sub-router)
- [Customizable response](https://github.com/gohf-http/gohf/blob/main/gohf_responses/json_response.go)
- Lightweight
- Based on [net/http](https://pkg.go.dev/net/http)

# 📍 Install GoHF

```sh
go get github.com/gohf-http/gohf
```

# 🪄 Features

### Feature: Easier Error Handing

```go
router.Handle("/greeting", func(c *gohf.Context) gohf.Response {
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
