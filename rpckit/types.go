package rpckit

import (
	"io"
	"net/http"
	"net/url"
)

type Params struct {
	Query url.Values `json:"query,omitempty"`
	Body  []byte     `json:"body,omitempty"`
}
type Request struct {
	Params
	Method     string `json:"method,omitempty"`
	RequestURI string `json:"requestURI,omitempty"`
}

func NewRequest(req *http.Request) (r *Request, err error) {
	r = &Request{
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
	Code        int         `json:"code,omitempty"`
	Msg         string      `json:"msg,omitempty"`
	Data        any         `json:"data,omitempty"`
	Rows        []any       `json:"rows,omitempty"`
	Total       int         `json:"total,omitempty"`
	ContentType ContentType `json:"contentType,omitempty"`
}

func NewResponse() *Response {
	return &Response{
		Rows: make([]any, 0),
	}
}
