# ExpressGo

## Introduction

As stated in the module name, this project aims to create a layer to use **Express.js** like API on top of Go standard **net/http**.

API in ExpressGo aligns to specifications of Express.js 5.x API [Reference](https://expressjs.com/en/5x/api.html).

ExpressGo leveraged **ServeMux** in **net/http** and would create a custom **ServeMux** with the following configurations.

### Default Configurations

1. **Host** is not allowed in path matching. Use `allowHost: true` in `App` to allow host in path matching.
2. All path matching is precise. Use `coarse: true` in `App` to fall back to default path matching of **ServeMux**.
3. Path matching is case insensitive. Use `app.use("case sensitive routing", true)` to fall back to default behaviors of **ServeMux**.

## TODO

### Jump to Next Route if Provided

1. redirect the request to another matching path in the handler

### Parse Params & Query String

1. parse params with the form from `:id` to `{id}`
2. set query string pairs into req.query

## Warning

This is currently still a hobby project for learning programming language Go. The module did not go through thorough testings and optimizations. Please use it at your own risk as stated in [License](./LICENSE).
