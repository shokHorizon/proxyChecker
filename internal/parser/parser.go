package parser

import (
	"context"
	"github.com/shokHorizon/proxyChecker/config"
	"github.com/shokHorizon/proxyChecker/internal/models"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type ProxyParser interface {
	GetUrls() []string
	GetName() string
	Parse(ctx context.Context, url string) ([]string, error)
}

type GenericParser struct {
	Re     regexp.Regexp
	Header http.Header
	Urls   []string
	Name   string
}

func NewFromConfig(cfg config.Source) *GenericParser {
	parser := &GenericParser{
		Re:     *regexp.MustCompile(cfg.Regexp),
		Header: http.Header{},
		Urls:   cfg.Urls,
		Name:   cfg.Name,
	}

	for _, header := range cfg.Headers {
		keyVal := strings.SplitN(header, ": ", 2)
		parser.Header.Set(keyVal[0], keyVal[1])
	}

	return parser
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

func (gp GenericParser) Parse(ctx context.Context, url string) ([]models.Proxy, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = gp.Header

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

	proxies := make([]models.Proxy, 0, 100)
	matches := gp.GetRe().FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		ip := match[1]
		port := match[2]
		proxies = append(proxies, *models.NewProxy(models.JoinIpPort(ip, port), gp.Name, "http"))
	}

	return proxies, nil
}
