package signer

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SignUrlMiddleware struct {
	RequestURI bool
	signer     ISigner
	FieldSign  bool
	Duration   time.Duration
}

func NewSignUrlMiddleware(conf Enpryt) *SignUrlMiddleware {
	var signerinstance ISigner
	if conf.Method == "md5" {
		signerinstance = NewMd5Signer(conf.Secret)
	} else {
		signerinstance = NewSha256Signer(conf.Secret)
	}
	return &SignUrlMiddleware{
		RequestURI: true,
		signer:     signerinstance,
		Duration:   time.Second * 3600,
	}
}
func (s *SignUrlMiddleware) Sign(requestURI string) (sign string, expireAt int64) {
	data := map[string]string{
		"requestURI": requestURI,
	}
	expireAt = time.Now().Add(s.Duration).Unix()
	sign, _ = s.signer.GenerateSignature(data, expireAt)
	return
}

// 鉴权
// https://tyest.turingdance.com/res/2025/06/21/178937897892738273.png?expireAt=105451244512&sign=09232832892838923923
func (s *SignUrlMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// 如果通过白名单
		// post:/a/{}
		requestURI := strings.TrimPrefix(req.URL.Path, "/")
		data := map[string]string{}
		//
		if s.RequestURI {
			data["requestURI"] = requestURI
		}
		//
		for key := range req.URL.Query() {
			data[key] = req.URL.Query().Get(key)
		}
		sign := data["sign"]
		expireAt, _ := strconv.ParseUint(data["expireAt"], 10, 64)
		delete(data, "sign")
		delete(data, "expireAt")
		//
		if ok, _ := s.signer.VerifySignature(data, sign, int64(expireAt)); ok {
			next.ServeHTTP(w, req)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	})
}
