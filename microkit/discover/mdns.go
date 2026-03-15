package discover

import (
	"context"
	"net"

	"github.com/grandcat/zeroconf"
)

type Mdns struct {
	app     string
	service string
	domain  string
	port    int
	text    []string
	ifaces  []net.Interface
}
type MdnsOption func(*Mdns)

func TCP() MdnsOption {
	return func(m *Mdns) {
		m.service = m.app + "._tcp"
	}
}
func App(app string) MdnsOption {
	return func(m *Mdns) {
		m.app = app
	}
}
func HTTP() MdnsOption {
	return func(m *Mdns) {
		m.service = m.app + "._http"
	}
}
func Type(tp string) MdnsOption {
	return func(m *Mdns) {
		m.text = append(m.text, "type="+tp)
	}
}
func Domain(domain string) MdnsOption {
	return func(m *Mdns) {
		m.domain = domain
	}
}
func Port[T int | int64 | int32 | uint | uint16 | uint32 | uint64](port T) MdnsOption {
	return func(m *Mdns) {
		m.port = int(port)
	}
}

func AddData(datamap ...string) MdnsOption {
	return func(m *Mdns) {
		m.text = append(m.text, datamap...)
	}
}

func NewMdns(app string, options ...MdnsOption) *Mdns {
	mdns := &Mdns{
		app:     app,
		service: app + "._tcp",
		domain:  "local.",
		ifaces:  nil,
		text:    []string{},
	}
	for _, v := range options {
		v(mdns)
	}
	return mdns
}

func (s *Mdns) Serve(options ...MdnsOption) (*zeroconf.Server, error) {
	for _, v := range options {
		v(s)
	}
	// 服务名、类型、域名、端口、本机IP
	server, err := zeroconf.Register(
		s.app,     // 服务名（唯一）
		s.service, // 服务类型（自定义）
		s.domain,  // 固定后缀
		s.port,    // 你的服务端口
		s.text,    // 自定义数据
		s.ifaces,
	)
	return server, err
}

func (s *Mdns) Browse(ctx context.Context, options ...MdnsOption) (entry chan<- *zeroconf.ServiceEntry, err error) {
	for _, v := range options {
		v(s)
	}
	resolver, _ := zeroconf.NewResolver(nil)
	entry = make(chan<- *zeroconf.ServiceEntry)
	err = resolver.Browse(ctx, s.service, s.domain, entry)
	return entry, err
}
