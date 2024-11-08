package restkit

import (
	"net/http"

	"github.com/turingdance/infra/wraper"
)

type (

	// HandlerFunc defines a function to serve HTTP requests.
	HandlerFunc func(ctx Context) (r *wraper.Response, err error)
	// MiddlewareFunc defines a function to process middleware.
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

type Route struct {
	Package     string
	Module      string
	Func        string
	Path        string
	Method      []string
	Comment     string
	HandlerFunc http.HandlerFunc
}
