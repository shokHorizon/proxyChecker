package provider

import "time"

type Proxy struct {
	IP       string
	Port     int
	Origin   string
	Protocol string
	FoundAt  time.Time
	Cooldown time.Time
	Retries  int
}
