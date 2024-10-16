# ExpressGo

## Introduction

As stated in the module name, this project aims to create a layer to use **Express.js** like API on top of Go standard **net/http**.

API in ExpressGo aligns to specifications of Express.js 5.x API [Reference](https://expressjs.com/en/5x/api.html).

ExpressGo leveraged **ServeMux** in **net/http** and would create a custom **ServeMux** with the following configurations.

### Default Configurations

1. **Host** is not allowed in path matching.
2. All path matching is precise.
3. Path matching is case insensitive.
4. Defining multiple lists of callbacks on the same route is allowed.

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

#### Add Callbacks to Routes

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

#### Path Params

Path params should be in the form of `:name`.

```go
app.Get("/user/:id", func(req *expressgo.Request, res *expressgo.Response, *expressgo.Next) {
    res.Send(req.Param["id"])
})

// Request: GET /user/101
// Respond: 101
```

To use separators, like hyphen (-) or dot (.):

```go
app.Get("/test/:one-:two-:three/:four.:five", func(req *expressgo.Request, res *expressgo.Response, *expressgo.Next) {
    lines := []string{}
    for k, v := range req.Params {
        lines = append(lines, fmt.Sprintf("%s: %s", k, v))
    }

    output := ""
    for _, line := range lines {
        output += line + "<br />"
    }
    res.Send(output)
})

// Request: GET /test/1-2-3/4.5
// Respond: one: 1<br />two: 2<br />three: 3<br />four: 4<br />five: 5<br />
```

Note:

1. Paths should not contain `{}`. ExpressGo would treat it as a literal and pass it down to `http.ServeMux`, and an error would occur.
2. Params should not have names ending with either `0H` or `0D`. These two strings are used for separators, including hyphens and dots.

#### Query String

WIP

### Response

WIP

### Next

At the current stage, it is still not possible to redifine function behaviors at runtime to mimic `next()` or `next('route')` usages in **Express.js**. Therefore, it is implemented this way to pass in a `*Next` pointer to a callback, so a callback could either use `next.Next = true` to activate the next callback or use `next.Route = true` to activate another list of callbacks defined on the same route. After the aforementioned `next.Next = true` or `next.Route = true` statement, remember to add `return` to exit the current callback if skipping any following logics is needed.

Note: `route` refers to the combination of `method` and `path`.

To run the next callback:

```go
// callback
func(*expressgo.Request, *expressgo.Response, next *expressgo.Next) {
    next.Next = true
    return
}
```

To run another list of callbacks defined on the same route:

```go
// callback
func(*expressgo.Request, *expressgo.Response, next *expressgo.Next) {
    next.Route = true
    return
}
```

Note: The next list refers to the list defined after the current list, in the order being called using the same `app.[Method]` on the same path.

## TODO

### Parse Query String

1. set query string pairs into req.query

### app.route()

1. chainable methods with path already included

### Router

1. mountable mini-app

## Warning

This is currently still a hobby project for learning programming language Go. The module did not go through thorough testings and optimizations. Please use it at your own risk as stated in [License](./LICENSE).
