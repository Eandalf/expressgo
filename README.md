# ExpressGo

## Introduction

As stated in the module name, this project aims to create a layer to use **Express.js** like API on top of Go standard **net/http**.

ExpressGo leveraged **ServeMux** in **net/http** and would create a custom **ServeMux** with the following configurations.

### Default Configurations

1. **Host** is not allowed in path matching. Use `allowHost: true` in `App` to allow host in path matching.
2. All path matching is precise. Use `coarse: true` in `App` to fall back to default path matching of **ServeMux**.
3. Path matching is case insensitive. Use `app.use("case sensitive routing", true)` to fall back to default behaviors of **ServeMux**.

## TODO

### Layer Between Express.js-like API to DefaultServeMux

1. Return an error if a host is trying to be registered into a path matching pattern.
2. Precise path matching using `{$}` at the end of every path matching pattern.
3. Case insensitive path matching, possible approach: [Kevin Gillette Re: [go-nuts] http.HandleFunc case insensitive path match, default match](https://groups.google.com/g/golang-nuts/c/M-_CyKCSGiA/m/-Z03K33HHRUJ).

## Warning

This is currently still a hobby project for learning programming language Go. The module did not go through thorough testings and optimizations. Please use it at your own risk as stated in [License](./LICENSE).
