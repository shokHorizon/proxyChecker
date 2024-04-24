package pinger

import (
	"context"
	"crypto/tls"
	"github.com/shokHorizon/proxyChecker/config"
	"github.com/shokHorizon/proxyChecker/internal/models"
	"golang.org/x/net/proxy"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"time"
)

const (
	checkURL2 = "https://2ip.ru"
	checkURL  = "https://steamcommunity.com/market/search/render/?query=&start=0&count=100&search_descriptions=0&sort_column=popular&sort_dir=desc&appid=730&norender=1"
	agent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type Pinger struct {
	cfg config.Pinger
}

func NewPinger(cfg config.Pinger) *Pinger {
	return &Pinger{
		cfg: cfg,
	}
}

func (p *Pinger) Run(ctx context.Context, r <-chan *models.Proxy, s chan<- *models.Proxy) {
	wg := errgroup.Group{}
	wg.SetLimit(p.cfg.Workers)
	for {
		select {
		case pr, ok := <-r:
			if !ok {
				break
			}
			wg.Go(
				func() error {
					err := CheckProxyHTTP(ctx, pr)
					if err != nil {
						err = CheckProxySocks(ctx, pr)
						if err != nil {
							return err
						}
					}
					s <- pr
					return nil
				})
		case <-ctx.Done():
			break
		}
	}
	wg.Wait()
}

func CheckProxyHTTP(ctx context.Context, p *models.Proxy) error {
	pUrl, err := url.Parse(p.Url())
	if err != nil {
		return err
	}

	transport := &http.Transport{Proxy: http.ProxyURL(pUrl)}
	client := &http.Client{Transport: transport}

	ctxT, _ := context.WithTimeout(ctx, time.Second*10)

	req, err := http.NewRequestWithContext(ctxT, "GET", checkURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", agent)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

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

	req, err := http.NewRequestWithContext(ctx, "GET", checkURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", agent)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

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

	req, err := http.NewRequestWithContext(ctxT, "GET", checkURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", agent)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
