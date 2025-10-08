package server

import (
	"fmt"
	"net/http"
	"strings"

	constant "github.com/Jamaceat/liquibase-versioning-app/constants"
	"github.com/go-chi/chi"
)

type Router interface {
	Use(middlewares ...func(http.Handler) http.Handler)
	With(middlewares ...func(http.Handler) http.Handler) Router
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Patch(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
	Options(pattern string, h http.HandlerFunc)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type engine struct {
	prefix      string
	chi         *chi.Mux
	middlewares []func(http.Handler) http.Handler
}

func (router *engine) Use(middlewares ...func(http.Handler) http.Handler) {
	router.chi.Use(middlewares...)
}

// With is used to add middlewares to a specific endpoint in the router
func (router *engine) With(middlewares ...func(http.Handler) http.Handler) Router {
	return &engine{
		prefix:      router.prefix,
		chi:         router.chi,
		middlewares: middlewares,
	}
}

// Get adds an http GET route to the router. The parameter pattern is the endpoint direction.
func (router *engine) Get(pattern string, handler http.HandlerFunc) {
	pattern = strings.TrimLeft(pattern, constant.Slash)
	router.chi.With(router.middlewares...).Get(fmt.Sprintf(constant.DoublePath, router.prefix, pattern), handler)
}

// Post adds an http POST route to the router. The parameter pattern is the endpoint direction.
func (router *engine) Post(pattern string, handler http.HandlerFunc) {
	pattern = strings.TrimLeft(pattern, constant.Slash)
	router.chi.With(router.middlewares...).Post(fmt.Sprintf(constant.DoublePath, router.prefix, pattern), handler)
}

// Put adds an http PUT route to the router. The parameter pattern is the endpoint direction.
func (router *engine) Put(pattern string, handler http.HandlerFunc) {
	pattern = strings.TrimLeft(pattern, constant.Slash)
	router.chi.With(router.middlewares...).Put(fmt.Sprintf(constant.DoublePath, router.prefix, pattern), handler)
}

// Patch adds an http PATCH route to the router. The parameter pattern is the endpoint direction.
func (router *engine) Patch(pattern string, handler http.HandlerFunc) {
	pattern = strings.TrimLeft(pattern, constant.Slash)
	router.chi.With(router.middlewares...).Patch(fmt.Sprintf(constant.DoublePath, router.prefix, pattern), handler)
}

// Delete adds an http DELETE route to the router. The parameter pattern is the endpoint direction.
func (router *engine) Delete(pattern string, handler http.HandlerFunc) {
	pattern = strings.TrimLeft(pattern, constant.Slash)
	router.chi.With(router.middlewares...).Delete(fmt.Sprintf(constant.DoublePath, router.prefix, pattern), handler)
}

// Options adds an http Options route to the router. The parameter pattern is the endpoint direction.
func (router *engine) Options(pattern string, handler http.HandlerFunc) {
	pattern = strings.TrimLeft(pattern, constant.Slash)
	router.chi.With(router.middlewares...).Options(fmt.Sprintf(constant.DoublePath, router.prefix, pattern), handler)
}

// ServeHTTP is the method needed to pass the router as http.Handler
func (router *engine) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	router.chi.ServeHTTP(rw, r)
}

// URLParam returns the url parameter key inside the request r
func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
