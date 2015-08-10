package ctxh

import (
	"net/http"

	"golang.org/x/net/context"
)

// ContextHandler defines a handler which receives a passed context.Context
// with the standard ResponseWriter and Request. ServeHTTP should write
// the reply headers and data to the ResponseWriter, cancel the ctx, and
// then return.
//
// Use a HandlerAdapter to wrap a ContextHandler as a http.Handler for
// compatability with ServeMux and middlewares. Middlewares which do not
// pass context.Context should come strictly before or after those that do,
// otherwise contexts will be lost.
type ContextHandler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

// ContextHandlerFunc type is an adapter to allow the use of an ordinary
// function as a ContextHandler. If f is a function with the correct
// signature, ContextHandlerFunc(f) is a ContextHandler that calls f.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTP calls the function f(ctx, w, req).
func (f ContextHandlerFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	f(ctx, w, req)
}

// HandlerAdapter wraps a ContextHandler to implement the http.Handler
// interface.
type HandlerAdapter struct {
	ctx     context.Context
	handler ContextHandler
}

// NewHandlerAdapter returns a http.Handler HandlerAdapter which wraps the
// given ContextHandler and creates a background context.Context.
func NewHandlerAdapter(handler ContextHandler) *HandlerAdapter {
	return &HandlerAdapter{
		ctx:     context.Background(),
		handler: handler,
	}
}

func (a *HandlerAdapter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.handler.ServeHTTP(a.ctx, w, req)
}
