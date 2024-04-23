package source

import (
	"github.com/shokHorizon/proxyChecker/internal/parser"
	"regexp"
)

func NewFreeProxyWorld() parser.GenericParser {
	re := regexp.MustCompile(`<tr>\n<td class="show-ip-div">\n(\d{1,4}.\d{1,4}.\d{1,4}.\d{1,4})\n</td>\n<td>\n<a.*?>(\d{1,5})</a>`)
	return parser.GenericParser{
		Re: *re,
		Urls: []string{
			"https://www.freeproxy.world/?type=http&anonymity=&country=&speed=&port=&page=1",
			"https://www.freeproxy.world/?type=http&anonymity=&country=&speed=&port=&page=2",
			"https://www.freeproxy.world/?type=http&anonymity=&country=&speed=&port=&page=3",
			"https://www.freeproxy.world/?type=http&anonymity=&country=&speed=&port=&page=4",
		},
		Name: "FreeProxyWorld",
	}
}
