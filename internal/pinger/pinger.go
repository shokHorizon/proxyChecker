package pinger

import (
	"context"
	"crypto/tls"
	"github.com/shokHorizon/proxyChecker/internal/models"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"time"
)

func CheckProxy(ctx context.Context, p *models.Proxy) error {
	pUrl, err := url.Parse(p.Url())
	if err != nil {
		return err
	}

	transport := &http.Transport{Proxy: http.ProxyURL(pUrl)}
	client := &http.Client{Transport: transport}

	ctxT, _ := context.WithTimeout(ctx, time.Second*10)

	req, err := http.NewRequestWithContext(ctxT, "GET", "https://steamcommunity.com/market/search/render/?query=&start=0&count=100&search_descriptions=0&sort_column=popular&sort_dir=desc&appid=730&norender=1", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func CheckProxySocks(ctx context.Context, p *models.Proxy) error {
	p.Protocol = "socks5"
	dialer, err := proxy.SOCKS5("tcp", p.Address, nil, proxy.Direct)
	if err != nil {
		return err
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: transport}

	ctx, _ = context.WithTimeout(ctx, time.Second*10)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://steamcommunity.com/market/search/render/?query=&start=0&count=100&search_descriptions=0&sort_column=popular&sort_dir=desc&appid=730&norender=1", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func CheckProxyHTTPS(ctx context.Context, p *models.Proxy) error {
	p.Protocol = "https"
	pUrl, err := url.Parse(p.Url())
	if err != nil {
		return err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(pUrl),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: transport}

	ctxT, _ := context.WithTimeout(ctx, time.Second*10)

	req, err := http.NewRequestWithContext(ctxT, "GET", "https://steamcommunity.com/market/search/render/?query=&start=0&count=100&search_descriptions=0&sort_column=popular&sort_dir=desc&appid=730&norender=1", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
