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

// handler wraps a ContextHandler to implement the http.Handler interface.
type handler struct {
	ctx     context.Context
	handler ContextHandler
}

// NewHandler returns an http.Handler which wraps the given ContextHandler
// and creates a background context.Context.
func NewHandler(h ContextHandler) http.Handler {
	return &handler{
		ctx:     context.Background(),
		handler: h,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handler.ServeHTTP(h.ctx, w, req)
}

// ContextHandlerFuncWithError is an adapter func to allow a context handler
// function which returns an AppError to be used as a ContextHandler. If f is
// a function with the correct signature, ContextHandlerFuncWithError(f) is a
// ContextHandler which calls f and writes the AppError if non-nil.
type ContextHandlerFuncWithError func(context.Context, http.ResponseWriter, *http.Request) *AppError

func (fn ContextHandlerFuncWithError) ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	if appErr := fn(ctx, w, req); appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
	}
}

// AppError bundles error data and should only be used as a return type for
// ContextHandlerFuncWithError functions.
type AppError struct {
	Error   error
	Message string
	Code    int
}
