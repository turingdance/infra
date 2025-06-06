package rpckit

import (
	"net/http"
)

type Gateway interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	Start()
	Stop()
}
