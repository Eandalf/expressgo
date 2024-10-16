# ExpressGo

## Introduction

As stated in the module name, this project aims to create a layer to use **Express.js** like API on top of Go standard **net/http**.

API in ExpressGo aligns to specifications of Express.js 5.x API [Reference](https://expressjs.com/en/5x/api.html).

ExpressGo leveraged **ServeMux** in **net/http** and would create a custom **ServeMux** with the following configurations.

### Default Configurations

1. **Host** is not allowed in path matching.
2. All path matching is precise.
3. Path matching is case insensitive.

To alter the behavior back to defaults of **net/http**:

```go
config := expressgo.Config{}
config.AllowHost = true // to allow host
config.Coarse = true // to opt-out precise path matching

app := expressgo.CreateServer(config)
app.Use("case sensitive routing", true) // to use case sensitive path matching
```

## Usage

### App

#### Create a Server

```go
config := expressgo.Config{} // optional
app := expressgo.CreateServer(config)
// or, without a config
// app := expressgo.CreateServer()
```

#### Add callbacks to routes

```go
// Get(string, func(*expressgo.Request, *expressgo.Response, *expressgo.Next))
app.Get("/hello", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    res.Send("Hello")
})

// With the style of middlewares
// Get(string, func1, func2, func3, ...)
```

#### Listen to a Port and Serve HTTP

```go
// Listen(int)
app.Listen(8080) // 8080 is the port number
```

### Request

WIP

### Response

WIP

### Next

WIP

## TODO

### Parse Params & Query String

1. parse params with the form from `:id` to `{id}`
2. set query string pairs into req.query

## Warning

This is currently still a hobby project for learning programming language Go. The module did not go through thorough testings and optimizations. Please use it at your own risk as stated in [License](./LICENSE).
