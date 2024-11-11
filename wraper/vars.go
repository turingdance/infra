package wraper

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func MuxIntVar[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](req *http.Request, key string, defaultvalue T) (result T, err error) {
	value, ok := mux.Vars(req)[key]
	if !ok {
		err = errors.New("变量不存在")
		return
	}
	iresult, err := strconv.Atoi(value)
	return T(iresult), err
}

func MuxFloatVar[T float32 | float64](req *http.Request, key string, defaultvalue T) (result T, err error) {
	value, ok := mux.Vars(req)[key]
	if !ok {
		err = errors.New("变量不存在")
		return
	}
	iresult, err := strconv.ParseFloat(value, 64)
	return T(iresult), err
}

func MuxStringVar[T string](req *http.Request, key string, defaultvalue T) (result T, err error) {
	value, ok := mux.Vars(req)[key]
	if !ok {
		err = errors.New("变量不存在")
		return
	}
	return T(value), err
}
