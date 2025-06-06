package rpcximpl

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rpcxio/libkv/store"
	rpcxredisclient "github.com/rpcxio/rpcx-redis/client"
	"github.com/smallnest/rpcx/client"
	"github.com/turingdance/infra/rpckit"
	"github.com/turingdance/infra/wraper"
)

type Rpcxapp struct {
	host     string
	port     int
	domain   string
	provider Provider
	server   *http.Server
}

func NewRpcxapp(domain string) *Rpcxapp {
	r := &Rpcxapp{
		host:   "",
		port:   8089,
		domain: domain,
	}
	return r
}
func (s *Rpcxapp) Port(port int) *Rpcxapp {
	s.port = port
	return s
}
func (s *Rpcxapp) Host(host string) *Rpcxapp {
	s.host = host
	return s
}
func (s *Rpcxapp) Provider(provider Provider) *Rpcxapp {
	s.provider = provider
	return s
}
func (g *Rpcxapp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//faas.turingdance.com/appname/module/action
	//domain := strings.Split(req.Host, ".")[0]
	arrs := strings.Split(req.URL.Path, "/")
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
			wraper.Error(err).Encode(w)
			return
		}
		xclient = client.NewXClient(module, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	}
	if g.provider.Type == MDNS {
		d, err := client.NewMDNSDiscovery(module, time.Second*5, time.Second*5, g.domain+"/"+appname)
		if err != nil {
			wraper.Error(err).Encode(w)
			return
		}
		xclient = client.NewXClient(module, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	}
	defer xclient.Close()

	request, err := rpckit.NewRequest(req)
	if err != nil {
		wraper.Error(err).Encode(w)
		return
	}
	reply := rpckit.NewResponse()
	err = xclient.Call(context.Background(), strings.ToUpper(action[:1])+action[1:], request, reply)
	if err != nil {
		wraper.Error(err).Encode(w)
		return
	} else {
		wraper.OkData(reply).Encode(w)
	}
}
func (g *Rpcxapp) Start() {
	addr := fmt.Sprintf("%s:%d", g.host, g.port)
	g.server = &http.Server{Addr: addr, Handler: g}
	println("run @", addr)
	err := g.server.ListenAndServe()
	if err != nil {
		println("catch error ", err.Error())
	}
}
func (g *Rpcxapp) Stop() {
	g.server.Shutdown(context.Background())
}
