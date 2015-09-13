# ctxh [![GoDoc](http://godoc.org/github.com/dghubble/ctxh?status.png)](http://godoc.org/github.com/dghubble/ctxh)

Package ctxh defines the common `ContextHandler` and associated adapters to avoid repeating it many of my packages. 

    type ContextHandler interface {
        ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
    }

The package will be marked deprecated if this (very common) interface becomes available in a standard library package.

## License

[Public Domain](LICENSE)



