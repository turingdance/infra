package ssekit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type DataFormat string

const (
	MarkDown DataFormat = "markdown"
	Html     DataFormat = "html"
	String   DataFormat = "string"
	JSON     DataFormat = "json"
)

func (s DataFormat) Description() string {
	ret := "字符串"
	switch s {
	case String:
		ret = "字符串"
	case MarkDown:
		ret = "MarkDown 文本"
	case Html:
		ret = "html网页"
	case JSON:
		ret = "Json格式"
	}
	return ret
}

type SseCtrl struct {
	current   int
	b64wraper bool
	writer    http.ResponseWriter
	format    DataFormat
	flusher   http.Flusher
}

func New(w http.ResponseWriter, opts ...SSEOption) *SseCtrl {
	flusher, _ := w.(http.Flusher)
	s := &SseCtrl{
		current:   1,
		format:    String,
		writer:    w,
		b64wraper: true,
		flusher:   flusher,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *SseCtrl) fromstring(data string) (frame *Frame, err error) {

	frame = NewFrame()
	if !s.b64wraper {
		frame.Data = data
	} else {
		frame.Data = base64.StdEncoding.EncodeToString([]byte(data))
	}
	return
}
func (s *SseCtrl) fromobject(data any) (frame *Frame, err error) {
	frame = NewFrame()
	bts, err := json.Marshal(data)
	if s.b64wraper {
		frame.Data = base64.StdEncoding.EncodeToString(bts)
	} else {
		frame.Data = string(bts)
	}
	return
}

func (s *SseCtrl) EncodeWith(c DataFormat) (frame *SseCtrl) {
	s.format = c
	return s
}

func (s *SseCtrl) WraperWithBase64(flag bool) (frame *SseCtrl) {
	s.b64wraper = flag
	return s
}

func (s *SseCtrl) WriteString(data string) (err error) {

	frame, err := s.fromstring(data)
	s.current = s.current + 1
	frame.Id = s.current
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(s.writer, "id:%d\ndata:%s\nevent:%s\n\n", frame.Id, frame.Data, frame.Event)
	if err != nil {
		return
	}
	s.flusher.Flush()
	return err
}

func (s *SseCtrl) WriteFrame(frame *Frame) (err error) {
	s.current = s.current + 1
	frame.Id = s.current
	_, err = fmt.Fprintf(s.writer, "id:%d\ndata:%s\nevent:%s\n\n", frame.Id, frame.Data, frame.Event)
	if err != nil {
		return
	}
	s.flusher.Flush()
	return err
}

func (s *SseCtrl) WriteObject(data any) (err error) {
	frame, err := s.fromobject(data)
	s.current = s.current + 1
	frame.Id = s.current
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(s.writer, "id:%d\ndata:%s\nevent:%s\n\n", frame.Id, frame.Data, frame.Event)
	if err != nil {
		return
	}
	s.flusher.Flush()
	return err
}
