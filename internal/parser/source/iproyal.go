package source

import (
	"github.com/shokHorizon/proxyChecker/internal/parser"
	"regexp"
)

func NewIpRoyal() parser.GenericParser {
	re := regexp.MustCompile(`<div class="flex items-center astro-lmapxigl">(\d{1,4}.\d{1,4}.\d{1,4}.\d{1,4})</div><div class="flex items-center astro-lmapxigl">(\d{1,5})</div>`)
	return parser.GenericParser{
		Re: *re,
		Urls: []string{
			"https://iproyal.com/free-proxy-list/?page=1&entries=100",
			"https://iproyal.com/free-proxy-list/?page=2&entries=100",
			"https://iproyal.com/free-proxy-list/?page=3&entries=100",
			"https://iproyal.com/free-proxy-list/?page=4&entries=100",
			"https://iproyal.com/free-proxy-list/?page=5&entries=100",
			"https://iproyal.com/free-proxy-list/?page=6&entries=100",
		},
		Name: "IpRoyal",
	}
}
