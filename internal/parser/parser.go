package parser

import (
	"context"
	"io"
	"net/http"
	"regexp"
)

type ProxyParser interface {
	GetUrls() []string
	GetName() string
	Parse(ctx context.Context, url string) ([]string, error)
}

type GenericParser struct {
	Re   regexp.Regexp
	Urls []string
	Name string
}

func (gp GenericParser) GetRe() *regexp.Regexp {
	return &gp.Re
}

func (gp GenericParser) GetUrls() []string {
	return gp.Urls
}

func (gp GenericParser) GetName() string {
	return gp.Name
}

func (gp GenericParser) Parse(ctx context.Context, url string) ([]string, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(bytes)

	addresses := make([]string, 0, 4)
	matches := gp.GetRe().FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		ip := match[1]
		port := match[2]
		addresses = append(addresses, ip+":"+port)
	}

	return addresses, nil
}
