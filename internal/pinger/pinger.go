package pinger

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"time"
)

func CheckProxy(ctx context.Context, proxyURL string) error {

	proxy, err := url.Parse("http://" + proxyURL)
	if err != nil {
		return err
	}

	transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
	client := &http.Client{Transport: transport}

	ctxT, _ := context.WithTimeout(ctx, time.Second*30)

	req, err := http.NewRequestWithContext(ctxT, "GET", "https://steamcommunity.com/", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return CheckProxyHTTPS(ctx, proxyURL)
	}

	return nil
}

func CheckProxySocks(ctx context.Context, proxyURL string) error {
	dialer, err := proxy.SOCKS5("tcp", "socks5://"+proxyURL, nil, proxy.Direct)
	if err != nil {
		return err
	}

	transport := &http.Transport{Dial: dialer.Dial}
	client := &http.Client{Transport: transport}

	ctx, _ = context.WithTimeout(ctx, time.Second*30)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://github.com/", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Socks found")
		return fmt.Errorf("proxy returned status %d", resp.StatusCode)
	}

	return nil
}

func CheckProxyHTTPS(ctx context.Context, proxyURL string) error {
	proxy, err := url.Parse("https://" + proxyURL)
	if err != nil {
		return err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: transport}

	ctxT, _ := context.WithTimeout(ctx, time.Second*30)

	req, err := http.NewRequestWithContext(ctxT, "GET", "https://github.com/", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return CheckProxySocks(ctx, proxyURL)
	}

	return nil
}
