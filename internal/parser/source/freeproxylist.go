package source

import (
	"github.com/shokHorizon/proxyChecker/internal/parser"
	"regexp"
)

func NewFreeProxyList() parser.GenericParser {
	re := regexp.MustCompile(`<tr><td>(\d{1,4}.\d{1,4}.\d{1,4}.\d{1,4})</td><td>(\d{1,5})</td>.*?</tr>`)
	return parser.GenericParser{
		Re: *re,
		Urls: []string{
			"https://free-proxy-list.net/",
		},
		Name: "FreeProxyList",
	}
}
