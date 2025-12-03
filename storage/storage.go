package storage

import (
	"net/http"
)

// oss://[bucket]/[key] => https://bucket.host/key
// local://[bucket]/[key] => https://host/buckey/key
type Response struct {
	Size         uint       `json:"size"`
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	Host         string     `json:"host"`
	Bucket       string     `json:"bucket"`
	Driver       DriverType `json:"driver"`
	SSL          bool       `json:"ssl"`
	AuthRequired bool       `json:"authRequired"`
}

// 这一层里面实现router
type Storage interface {
	Upload(r *http.Request) Response
}
