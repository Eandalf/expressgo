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
app.Set("case sensitive routing", true) // to use case sensitive path matching
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

Path params should be in the form of `:name`. A valid param name has the form of `[A-Za-z_][A-Za-z0-9_]*`, starting with A-Z, a-z, or underscore (\_), and concatenated with A-Z, a-z, 0-9, or underscore (\_).

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

> Note:
>
> 1. Paths should not contain `{}`. ExpressGo would treat it as a literal and pass it down to `http.ServeMux`, and an error would occur.
> 2. Params should not have names ending with either `0H` or `0D`. These two strings are used for separators, including hyphens and dots.

#### Query String

Query string could be read from `req.Query["key"]`.

```go
app.Get("/test/query", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    res.Send(req.Query["id"])
})

// Request: GET /test/query?id=101
// Respond: 101
```

> Note:
>
> 1. Query string would be parsed no matter with which http method.
> 2. Only the first value of a key from the query string is parsed.

#### Body (JSON)

**ExpressGo** provides a package under [github.com/Eandalf/expressgo/bodyparser](https://github.com/Eandalf/expressgo/bodyparser) for parsing the body of a request.

`bodyparser.Json()` returns a parser as a middleware to parse received body stream with a specified type into `req.Body`. It defaults to use `expressgo.BodyJsonBase`, which is basically `map[string]interface{}`, as the received JSON type. Custom types could be supplied to the parser through `bodyparser.Json(bodyparser.JsonConfig{Receiver: &Test{}})` where `Test` is the name of the custom type. It is recommended to pass the pointer of the custom struct to `Receiver` option since the underlying decoder is `json.NewDecoder(...).Decode(...)` from **encoding/json**.

The parser leverages **encoding/json**. Hence, the custom struct should follow tag notations used in **encoding/json**.

For example,

```go
type Test struct {
    Test string `json:"test"`
}
```

To parse JSON with the default struct `expressgo.BodyJsonBase`:

```go
app.Post("/test/body/base", bodyparser.Json(), func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    if j, ok := req.Body.(expressgo.BodyBase); ok {
        if t, ok := j["test"]; ok {
            if s, ok := t.(string); ok {
                res.Send(s)
            }
        }
    }

    res.Send("body parsing failed")
})

// Request: POST /test/body/base/
// Body: '{"test":"test_string"}'
// Respond: test_string
```

To parse JSON with a custom struct:

```go
type Test struct {
    Test string `json:"test"`
}

app.Post("/test/body/type", bodyparser.Json(bodyparser.JsonConfig{Receiver: &Test{}}), func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    if t, ok := req.Body.(*Test); ok {
        res.Send(t.Test)
    }

    res.Send("body parsing failed")
})

// Request: POST /test/body/type/
// Body: '{"test":"test_string"}'
// Respond: test_string
```

> Note:
>
> 1. `req.Body` is typed as `interface{}`.
> 2. Although it is common to set `bodyParser.json()` as a global middleware in **Express.js**, with static type constraints in Go, it is not idiomatic to do so. Since it is common to have callbacks for POST requests expecting different DTOs, it is more suitable to place the JSON parser on each route as shown in the examples above.
> 3. `bodyparser.Json()` could not be invoked twice on the same route (same method and same path), the parser would consume the body stream, which would lead to nothing left for the coming parser to process. If two JSON parsers are invoked, the second one would be a no-op instead of raising the `io.EOF` error to the next error-handling callback.

### Response

WIP

### Next

At the current stage, it is still not possible to redifine function behaviors at runtime to mimic `next()` or `next('route')` usages in **Express.js**. Therefore, it is implemented this way to pass in a `*Next` pointer to a callback, so a callback could either use `next.Next = true` to activate the next callback or use `next.Route = true` to activate another list of callbacks defined on the same route. After the aforementioned `next.Next = true` or `next.Route = true` statement, remember to add `return` to exit the current callback if skipping any following logics is needed.

> Note:
>
> 1. `route` refers to the combination of `method` and `path`.
> 2. `next.Next` always takes precedence over `next.Route` if both are set.

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

> Note: The next list refers to the list defined after the current list, in the order being called using the same `app.[Method]` on the same path.

## Method

app.[Method]

### app.UseGlobal

To mount callbacks as middlewares to all paths with all http methods.

The order of invocation matters. The callbacks of `app.[Method]` defined before `app.UseGlobal` would be executed before the inserted middlewares using `app.UseGlobal`.

```go
app.UseGlobal(func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    req.Params["global"] = "global"
    // next.Route is recommended to be set to `true`, otherwise, nothing after the middleware could be executed
    next.Route = true
})

app.Get("/test/use/global", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    res.Send(req.Params["global"])
})

// Request: GET /test/use/global
// Respond: global
```

### app.Use

To mount callbacks as middlewares to the path with all http methods.

The order of invocation matters. The callbacks of `app.[Method]` defined before `app.Use` would be executed before the inserted middlewares using `app.Use`.

```go
app.Use("/test/use", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    req.Params["id"] = "101"
    // next.Route is recommended to be set to `true`, otherwise, nothing after the middleware could be executed
    next.Route = true
})

app.Get("/test/use", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    res.Send(req.Params["id"])
})

// Request: GET /test/use
// Respond: 101
```

### app.Get

For GET requests.

```go
app.Get("/", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    res.Send("Hello from root")
})

// Request: GET /
// Respond: Hello from root
```

### app.Post

For POST requests.

```go
app.Post("/test/body/base", bodyparser.Json(), func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    if j, ok := req.Body.(expressgo.BodyBase); ok {
        if t, ok := j["test"]; ok {
            if s, ok := t.(string); ok {
                res.Send(s)
            }
        }
    }

    res.Send("body parsing failed")
})

// Request: POST /test/body/base/
// Body: '{"test":"test_string"}'
// Respond: test_string
```

> Note:
>
> 1. Requests from http clients to POST paths need to have the path *very* precise. For example, `app.Post("/test/body/base", ...)` would need the path to be set to `/test/body/base/` in client requests.
> 2. This is caused by the default behavior of **ExpressGo** to make path precise and the default redirect http status code (301) used by **net/http**.
> 3. While making the path precise, **ExpressGo** actually forces each path to have a trailing slash (/).
> 4. While an http client sends a request to the originally designated path (`/path`), **net/http** would send a redirect with status code 301 to point to `/path/`.
> 5. This would cause the client to drop the request body and resend the request through GET method as per status code 301 indicated.
> 6. Related issue: [golang/go#60769](https://github.com/golang/go/issues/60769)

## Error Handling

If any error is intended to be handled by other callbacks, set `next.Error = error` to pass the error to any error handler behind.

After an error handler is triggered, the error is seemed as consumed. If the error needs to be passed to another error handler, set `next.Error = error` in the current error handler to pass the error down to the next error handler.

Error handlers are set with similar logics as `app.Use` and `app.UseGlobal`, so the order of invocation matters.

`app.UseError` and `app.UseGlobalError` are often used at the very end of all `app.[Method]` calls.

### app.UseError

To mount an error handler on a path with all http methods.

```go
app.Get("/test/error", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    next.Err = errors.New("raised error in /test/error")
    return // optional, to skip any logics behind
})

app.UseError("/test/error", func(err error, req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    res.Send(err.Error())
})

// Request: GET /test/error
// Respond: raised error in /test/error
```

### app.UseGlobalError

To mount an error handler to all routes.

```go
app.Get("/test/error/1", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    next.Err = errors.New("raised error in /test/error/1")
    return // optional, to skip any logics behind
})

app.Get("/test/error/2", func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    next.Err = errors.New("raised error in /test/error/2")
    return // optional, to skip any logics behind
})

app.UseGlobalError(func(err error, req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
    res.Send(err.Error())
})

// Request: GET /test/error/1
// Respond: raised error in /test/error/1

// Request: GET /test/error/2
// Respond: raised error in /test/error/2
```

## TODO

### app.route()

1. chainable methods with path already included

### Router

1. mountable mini-app

## Warning

This is currently still a hobby project for learning programming language Go. The module did not go through thorough testings and optimizations. Please use it at your own risk as stated in [License](./LICENSE).
