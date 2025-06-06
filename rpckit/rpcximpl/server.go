package rpcximpl

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/libkv/store"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"github.com/turingdance/infra/netkit"
)

type RegisterProviderType string

const REDIS RegisterProviderType = "redis"
const CONSUL RegisterProviderType = "consul"
const MDNS RegisterProviderType = "mdns"

type Provider struct {
	Type           RegisterProviderType `json:"type"`
	ServerURL      string               `json:"serverURL"`
	UpdateInterval time.Duration        `json:"updateInterval"`
	Password       string               `json:"password"`
	Username       string               `json:"username"`
	Dbname         string               `json:"dbname"`
}
type Server struct {
	network   string
	host      string
	port      int
	domain    string
	name      string
	provider  Provider
	rpcserver *server.Server
}
type ServerOption func(*Server)

func WithAddr(addr string) ServerOption {
	return func(s *Server) {
		arr := strings.Split(addr, ":")
		s.host = arr[0]
		s.port, _ = strconv.Atoi(arr[1])
	}
}
func WithName(name string) ServerOption {
	return func(s *Server) {
		s.name = name
	}
}
func RandomPort() ServerOption {
	return func(s *Server) {
		s.port, _ = netkit.RandomPort()
	}
}
func WithProvider(p Provider) ServerOption {
	return func(s *Server) {
		s.provider = p
	}
}

func NewServer(domain, name string, options ...ServerOption) *Server {
	r := &Server{
		rpcserver: server.NewServer(),
		network:   "tcp",
		host:      "",
		port:      8089,
		name:      name,
		domain:    domain,
	}
	for _, o := range options {
		o(r)
	}
	return r
}
func (s *Server) Provider(provider Provider) *Server {
	s.provider = provider
	s.setupregister()
	return s
}
func (s *Server) Serve() *Server {
	s.rpcserver.Serve(s.network, fmt.Sprintf("%s:%d", s.host, s.port))
	return s
}
func (s *Server) Stop() {

}
func (s *Server) Register(rcvr interface{}, metadata string) *Server {
	s.rpcserver.Register(rcvr, metadata)
	return s
}
func (s *Server) RegisterName(name string, rcvr interface{}, metadata string) *Server {
	s.rpcserver.RegisterName(name, rcvr, metadata)
	return s
}
func (s *Server) RegisterFunction(servicePath string, fn interface{}, metadata string) *Server {
	s.rpcserver.RegisterFunction(servicePath, fn, metadata)
	return s
}
func (s *Server) RegisterFunctionName(servicePath string, name string, fn interface{}, metadata string) *Server {
	s.rpcserver.RegisterFunctionName(servicePath, name, fn, metadata)
	return s
}

func (s *Server) setupregister() (err error) {
	serviceAddress := fmt.Sprintf("%s@%s:%d", s.network, s.host, s.port)
	metri := metrics.DefaultRegistry
	domain := s.name
	if s.provider.Type == MDNS {
		r := serverplugin.NewMDNSRegisterPlugin(serviceAddress, s.port, metri, s.provider.UpdateInterval, domain)
		err = r.Start()
		s.rpcserver.Plugins.Add(r)
	}
	if s.provider.Type == REDIS {
		var option *store.Config = &store.Config{}
		if s.provider.Password != "" {
			option.Password = s.provider.Password
		}
		if s.provider.Username != "" {
			option.Username = s.provider.Username
		}
		if s.provider.Dbname != "" {
			option.Bucket = s.provider.Dbname
		}
		r := &serverplugin.RedisRegisterPlugin{
			ServiceAddress: serviceAddress,
			RedisServers:   []string{s.provider.ServerURL},
			BasePath:       s.domain + "/" + s.name,
			UpdateInterval: s.provider.UpdateInterval,
			Options:        option,
		}
		s.rpcserver.Plugins.Add(r)
	}
	return err
}
