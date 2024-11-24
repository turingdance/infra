package restkit

import (
	"net/http"

	"github.com/turingdance/infra/wraper"
)

type (

	// HandlerFunc defines a function to serve HTTP requests.
	HandlerFunc func(ctx Context) (r *wraper.Response, err error)
	// MiddlewareFunc defines a function to process middleware.
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

type Route struct {
	Package     string           //包名称
	Module      string           //模块抹茶
	Func        string           //函数名称
	Path        string           //路径
	Method      []string         //方法
	Comment     string           // 配置
	App         string           //所属应用
	HandlerFunc http.HandlerFunc //对应处理器
}
