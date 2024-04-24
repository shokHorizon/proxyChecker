package provider

import (
	"context"
	"fmt"
	"github.com/shokHorizon/proxyChecker/config"
	"github.com/shokHorizon/proxyChecker/internal/models"
	"time"
)

type Provider struct {
	config.Provider
	pool map[string]*models.Proxy
}

func NewProvider(cfg config.Provider) *Provider {
	return &Provider{
		Provider: cfg,
		pool:     make(map[string]*models.Proxy),
	}
}

func (p *Provider) Run(ctx context.Context, receiver <-chan *models.Proxy) {
	for {
		select {
		case pr, ok := <-receiver:
			if !ok {
				return
			}
			_, ok = p.pool[pr.Address]
			if ok {
				fmt.Println("Duplicate proxy: ", pr.Address)
				continue
			}
			p.pool[pr.Address] = pr
			fmt.Printf("%.3d %20.20s %28.28s %s\n", len(p.pool), pr.Origin, pr.Address, pr.Protocol)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Provider) Get() *models.Proxy {
	t := time.Now()
	for _, pr := range p.pool {
		if pr.Cooldown.Before(t) {
			pr.Cooldown = pr.Cooldown.Add(p.BookTime)
			return pr
		}
	}
	return nil
}

func (p *Provider) Free(pr *models.Proxy) {
	proxy, ok := p.pool[pr.Address]
	if !ok {
		return
	}
	proxy.Cooldown = time.Now()
}

func (p *Provider) Dead(pr *models.Proxy) {
	proxy, ok := p.pool[pr.Address]
	if !ok {
		return
	}
	proxy.Retries++
	proxy.Cooldown = time.Now().Add(p.CoolTime)
}
