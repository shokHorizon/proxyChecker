package source

import (
	"github.com/shokHorizon/proxyChecker/internal/parser"
	"regexp"
)

func NewHideMyName() parser.GenericParser {
	re := regexp.MustCompile(`<tr><td>(\d{1,4}.\d{1,4}.\d{1,4}.\d{1,4})</td><td>(\d{1,5})</td>.*?</tr>`)
	return parser.GenericParser{
		Re: *re,
		Urls: []string{
			"https://hidemy.io/ru/proxy-list/countries/russian-federation/",
			"https://hidemy.io/ru/proxy-list/countries/germany/",
			"https://hidemy.io/ru/proxy-list/countries/netherlands/",
			"https://hidemy.io/ru/proxy-list/countries/ukraine/",
			"https://hidemy.io/ru/proxy-list/countries/france/",
			"https://hidemy.io/ru/proxy-list/countries/romania/",
		},
		Name: "HideMyName",
	}
}
