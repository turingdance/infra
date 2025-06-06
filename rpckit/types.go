package rpckit

import (
	"io"
	"net/http"
	"net/url"
)

type Params struct {
	Query url.Values `json:"query"`
	Body  []byte     `json:"body"`
}
type Request struct {
	Params
	Method     string `json:"method"`
	RequestURI string `json:"requestURI"`
}

func NewRequest(req *http.Request) (r Request, err error) {
	r = Request{
		Params: Params{
			Query: req.URL.Query(),
		},
		Method:     req.Method,
		RequestURI: req.RequestURI,
	}
	if req.Method == http.MethodPost || req.Method == http.MethodPut {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return r, err
		}
		r.Params.Body = append([]byte{}, body...)
	}
	return
}

type ContentType string
type Response struct {
	Code        int         `json:"code"`
	Msg         string      `json:"msg"`
	Data        any         `json:"data"`
	Rows        []any       `json:"rows"`
	Total       int         `json:"total"`
	ContentType ContentType `json:"contentType"`
}

func NewResponse() Response {
	return Response{
		Rows: make([]any, 0),
	}
}
