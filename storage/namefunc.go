package storage

import (
	"fmt"
	"mime/multipart"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/turingdance/infra/slicekit"
)

var nameFuncMap *sync.Map = &sync.Map{}
var currentid int64 = time.Now().Unix()

type Context struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
	Ext        string //fileext 后缀
	UserId     any
	TeamId     any
	TenantId   any
}
type NameFuncType func(ctx Context) string

func YYYY(ctx Context) string {
	return fmt.Sprintf("%04d", time.Now().Year())
}

func MM(ctx Context) string {
	return fmt.Sprintf("%02d", time.Now().Month())
}

func DD(ctx Context) string {
	return fmt.Sprintf("%02d", time.Now().Day())
}

func UUID(ctx Context) string {
	u := uuid.New()
	return u.String()
}
func TIMESTAMP(ctx Context) string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
func DATETIME(ctx Context) string {
	return fmt.Sprintf("%s%d", time.Now().Format("20060102150405"), atomic.AddInt64(&currentid, 1))
}
func FILENAME(ctx Context) string {
	return strings.Split(ctx.FileHeader.Filename, ".")[0]
}
func USERID(ctx Context) string {
	switch t := ctx.UserId.(type) {
	case string:
		return t
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", t)
	case float32, float64:
		return fmt.Sprintf("%f", t)
	default:
		return ""
	}

}
func TEAMID(ctx Context) string {
	switch t := ctx.TeamId.(type) {
	case string:
		return t
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", t)
	case float32, float64:
		return fmt.Sprintf("%f", t)
	default:
		return ""
	}
}
func TENANTID(ctx Context) string {
	switch t := ctx.TenantId.(type) {
	case string:
		return t
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", t)
	case float32, float64:
		return fmt.Sprintf("%f", t)
	default:
		return ""
	}
}
func Run(ctx Context, funs ...NameFuncType) string {
	tmp := []string{}
	for _, fun := range funs {
		tmp = append(tmp, fun(ctx))
	}

	return strings.Join(slicekit.Filter[string](tmp, func(item string, index int, slice []string) bool {
		return item != ""
	}), "/")
}
func Invoke(ctx Context, funs ...string) string {
	tmp := []string{}
	for _, funname := range funs {
		funname = strings.ToUpper(funname)
		fun, ok := nameFuncMap.Load(funname)
		if ok {
			handler := fun.(NameFuncType)
			tmp = append(tmp, handler(ctx))
		}
	}

	return strings.Join(slicekit.Filter[string](tmp, func(item string, index int, slice []string) bool {
		return item != ""
	}), "/")
}

func init() {
	nameFuncMap.Store("YYYY", NameFuncType(YYYY))
	nameFuncMap.Store("MM", NameFuncType(MM))
	nameFuncMap.Store("DD", NameFuncType(DD))
	nameFuncMap.Store("UUID", NameFuncType(UUID))
	nameFuncMap.Store("TIMESTAMP", NameFuncType(TIMESTAMP))
	nameFuncMap.Store("DATETIME", NameFuncType(DATETIME))
	nameFuncMap.Store("FILENAME", NameFuncType(FILENAME))
	nameFuncMap.Store("USERID", NameFuncType(USERID))
	nameFuncMap.Store("TEAMID", NameFuncType(TEAMID))
	nameFuncMap.Store("TENANTID", NameFuncType(TENANTID))
}
