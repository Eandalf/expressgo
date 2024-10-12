# ExpressGo

## Introduction

As stated in the module name, this project aims to create a layer to use **Express.js** like API on top of Go standard **net/http**.

## TODO

### To Register

1. check if passed in path pattern string is blank or not
2. check if handler is nil or not
3. check if the handler match the signature of ServeHTTP (we might be able to extend the definition of ServeHTTP)
4. parse the pattern string
5. check if the patterns conflict, the default behavior should be able to overwrite the existing path and handler, this could be an option in App object
6. if a conflict is detected, show the conflicting function locations (file & line), could use runtime.Caller, in backlog
7. goroutine safe (mux.mu.Lock()\n defer mux.mu.Unlock())

### To Serve

1. path matching to find the handler
2. pass the request to the handler

## Warning

This is currently still a hobby project for learning programming language Go. The module did not go through thorough testings and optimizations. Please use it at your own risk as stated in [License](./LICENSE).
