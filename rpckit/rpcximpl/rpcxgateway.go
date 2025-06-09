package rpcximpl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rpcxio/libkv/store"
	rpcxredisclient "github.com/rpcxio/rpcx-redis/client"
	"github.com/smallnest/rpcx/client"
	"github.com/turingdance/infra/logger"
	"github.com/turingdance/infra/rpckit"
	"github.com/turingdance/infra/wraper"
)

type Route struct {
	patern     string
	hander     http.Handler
	method     []string
	middleware []mux.MiddlewareFunc
}

func NewRoute(pater string, hander http.Handler) *Route {
	return &Route{
		patern: pater,
		hander: hander,
		method: []string{
			"post",
		},
		middleware: []mux.MiddlewareFunc{},
	}
}

func (r *Route) Patern(p string) *Route {
	r.patern = p
	return r
}
func (r *Route) Handle(h http.Handler) *Route {
	r.hander = h
	return r
}
func (r *Route) Method(m ...string) *Route {
	r.method = append(r.method, m...)
	return r
}
func (r *Route) Middleware(m ...mux.MiddlewareFunc) *Route {
	r.middleware = append(r.middleware, m...)
	return r
}

type Rpcxapp struct {
	host     string
	port     int
	domain   string
	provider Provider
	server   *http.Server
	rpcroute *Route
	routes   []*Route
	logger   logger.ILogger
	// 全局中间件
	middleware []mux.MiddlewareFunc
}

func NewRpcxapp(domain string) *Rpcxapp {
	r := &Rpcxapp{
		host:   "",
		port:   8089,
		domain: domain,
		rpcroute: &Route{
			patern:     "/{appname}/{module}/{action}",
			method:     []string{"post"},
			middleware: []mux.MiddlewareFunc{},
		},
		logger:     logger.Std(),
		routes:     make([]*Route, 0),
		middleware: []mux.MiddlewareFunc{},
	}
	return r
}
func (s *Rpcxapp) Port(port int) *Rpcxapp {
	s.port = port
	return s
}
func (s *Rpcxapp) UseLogger(l logger.ILogger) *Rpcxapp {
	s.logger = l
	return s
}
func (s *Rpcxapp) Host(host string) *Rpcxapp {
	s.host = host
	return s
}
func (s *Rpcxapp) AddRoute(route ...*Route) *Rpcxapp {
	s.routes = append(s.routes, route...)
	return s
}
func (s *Rpcxapp) Router(route *Route) *Rpcxapp {
	s.routes = append(s.routes, route)
	return s
}
func (s *Rpcxapp) Provider(provider Provider) *Rpcxapp {
	s.provider = provider
	return s
}
func (s *Rpcxapp) Use(mdw ...mux.MiddlewareFunc) *Rpcxapp {
	s.middleware = append(s.middleware, mdw...)
	return s
}

func (s *Rpcxapp) Ctxpath(path string) *Rpcxapp {
	s.rpcroute.patern = path
	return s
}
func (s *Rpcxapp) Method(method ...string) *Rpcxapp {
	s.rpcroute.method = method
	return s
}
func (g *Rpcxapp) NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	wraper.Error(errors.New("not found")).Encode(w)
}
func (g *Rpcxapp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	arrs := strings.Split(strings.TrimLeft(req.URL.Path, "/"), "/")
	g.logger.Debugf("[%s],[%s],%s", req.Method, req.RemoteAddr, req.RequestURI)
	if len(arrs) < 3 {
		g.logger.Errorf("[%s],[%s],%s,path is unavaliable", req.Method, req.RemoteAddr, req.RequestURI)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	appname := arrs[len(arrs)-3]
	module := strings.ToUpper(arrs[len(arrs)-2][:1]) + arrs[len(arrs)-2][1:]
	action := arrs[len(arrs)-1]
	var xclient client.XClient
	if g.provider.Type == REDIS {
		var option *store.Config = &store.Config{}
		if g.provider.Password != "" {
			option.Password = g.provider.Password
		}
		if g.provider.Username != "" {
			option.Username = g.provider.Username
		}
		if g.provider.Dbname != "" {
			option.Bucket = g.provider.Dbname
		}
		d, err := rpcxredisclient.NewRedisDiscovery((g.domain + "/" + appname), module, []string{g.provider.ServerURL}, option)
		if err != nil {
			g.logger.Errorf("[%s],[%s],%s,err=%s", req.Method, req.RemoteAddr, req.RequestURI, err.Error())
			wraper.Error(errors.New("服务未找到")).Encode(w)
			return
		}
		xclient = client.NewXClient(module, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	}
	if g.provider.Type == MDNS {
		d, err := client.NewMDNSDiscovery(module, time.Second*5, time.Second*5, g.domain+"/"+appname)
		if err != nil {
			g.logger.Errorf("[%s],[%s],%s,err=%s", req.Method, req.RemoteAddr, req.RequestURI, err.Error())
			wraper.Error(errors.New("服务未找到")).Encode(w)
			return
		}
		xclient = client.NewXClient(module, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	}
	defer xclient.Close()

	request, err := rpckit.NewRequest(req)
	if err != nil {
		g.logger.Errorf("[%s],[%s],%s,err=%s", req.Method, req.RemoteAddr, req.RequestURI, err.Error())
		wraper.Error(errors.New("参数解析有误")).Encode(w)
		return
	}
	reply := rpckit.NewResponse()
	err = xclient.Call(context.Background(), strings.ToUpper(action[:1])+action[1:], request, reply)
	if err != nil {
		g.logger.Errorf("[%s],[%s],%s,err=%s", req.Method, req.RemoteAddr, req.RequestURI, err.Error())
		wraper.Error(errors.New("系统内部发生错误")).Encode(w)
		return
	} else {
		wraper.OkData(reply).Encode(w)
	}
}
func (g *Rpcxapp) Start() {
	addr := fmt.Sprintf("%s:%d", g.host, g.port)
	g.rpcroute.hander = g
	routes := append([]*Route{}, g.routes...)
	router := mux.NewRouter()
	// 全局
	router.NotFoundHandler = g
	router.Use(g.middleware...)

	// 处理全部
	for _, v := range routes {
		// subrouter := router.NewRoute().Subrouter()
		// subrouter.Handle(v.patern, v.hander).Methods(v.method...)
		// subrouter.Use(v.middleware...)
		router.PathPrefix(v.patern).Handler(v.hander).Methods(v.method...)
		g.logger.Infof("register %s", v.patern)
	}
	g.server = &http.Server{Addr: addr, Handler: router}
	g.logger.Infof("micro run @ %s", addr)
	err := g.server.ListenAndServe()
	if err != nil {
		g.logger.Errorf("micro error ", err.Error())
	}
}
func (g *Rpcxapp) Stop() {
	g.server.Shutdown(context.Background())
	g.logger.Infof("micro stop ")
}
