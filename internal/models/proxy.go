package models

import (
	"fmt"
	"time"
)

type Proxy struct {
	Address  string
	Origin   string
	Protocol string
	FoundAt  time.Time
	Cooldown time.Time
	Retries  int
}

func NewProxy(address string, origin string, protocol string) *Proxy {
	return &Proxy{
		Address:  address,
		Origin:   origin,
		Protocol: protocol,
		FoundAt:  time.Now(),
		Cooldown: time.Now(),
		Retries:  0,
	}
}

func (p *Proxy) Url() string {
	return p.Protocol + "://" + p.Address
}

func JoinIpPort(ip, port string) string {
	return fmt.Sprintf("%s:%s", ip, port)
}
