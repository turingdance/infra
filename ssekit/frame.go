package ssekit

import (
	"encoding/base64"
	"encoding/json"
)

type FrameOption func(*Frame)
type Frame struct {
	Id    int    `json:"id"`
	Data  string `json:"data"`
	Event string `json:"event"`
}

func NewFrame(opts ...FrameOption) *Frame {
	r := &Frame{
		Id:    0,
		Event: "message",
	}
	for _, v := range opts {
		v(r)
	}
	return r
}

// WithFrameId 设置初始Frame的Id
func FrameId(id int) FrameOption {
	return func(s *Frame) {
		s.Id = id
	}
}

// WithFrameData 设置初始Frame的Data内容
func Data(data string) FrameOption {
	return func(s *Frame) {
		s.Data = data
	}
}

// WithFrameData 设置初始Frame的Data内容
func EncodeTOBase64(data string) FrameOption {
	return func(s *Frame) {
		s.Data = base64.StdEncoding.EncodeToString([]byte(data))
	}
}

// WithFrameData 设置初始Frame的Data内容
func EncodeTOJsonThenBase64(data any) FrameOption {
	return func(s *Frame) {
		ret, _ := json.Marshal(data)
		s.Data = base64.StdEncoding.EncodeToString(ret)
	}
}

// WithFrameEvent 设置初始Frame的Event类型（如"message"/"update"）
func Event(event string) FrameOption {
	return func(s *Frame) {
		s.Event = event
	}
}
