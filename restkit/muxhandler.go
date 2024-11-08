package restkit

import (
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

type MuxHandler struct {
	router *mux.Router
}

func NewMuxHandler() *MuxHandler {
	r := &MuxHandler{
		router: mux.NewRouter(),
	}
	return r
}
func (r *MuxHandler) Router(patern string, routes ...Route) *MuxHandler {
	router := r.router.PathPrefix(path.Join("/", patern)).Subrouter()
	for _, route := range routes {
		if !strings.Contains(route.Path, "{") {
			router.HandleFunc(route.Path, route.HandlerFunc).Methods(route.Method...)
		}
	}
	// 特殊字符的后面处理，基本上是pathvariable
	for _, route := range routes {
		if strings.Contains(route.Path, "{") {
			router.HandleFunc(route.Path, route.HandlerFunc).Methods(route.Method...)
		}
	}
	return r
}
func (s *MuxHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}
