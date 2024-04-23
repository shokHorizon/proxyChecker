package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/shokHorizon/proxyChecker/internal/parser"
	"github.com/shokHorizon/proxyChecker/internal/parser/source"
	"github.com/shokHorizon/proxyChecker/internal/pinger"
	"github.com/shokHorizon/proxyChecker/internal/saver"
	"sync"
	"time"
)

const (
	maxConn     = 5000
	maxWorkers  = 500
	maxRequests = 1
	steamPages  = 200
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	chProxy := make(chan string, 400)
	chProved := make(chan string, 400)
	pMutex := sync.Mutex{}
	proxies := 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		pwg := sync.WaitGroup{}
		parsers := []parser.ProxyParser{
			source.NewFreeProxyList(),
			source.NewFreeProxyWorld(),
			source.NewHideMyLife(),
			source.NewHideMyName(),
			source.NewIpRoyal(),
		}
		for _, page := range parsers {
			pwg.Add(1)
			page := page
			go func() {
				defer pwg.Done()
				for _, url := range page.GetUrls() {
					ctx, _ := context.WithTimeout(ctx, time.Second*5)
					result, err := page.Parse(ctx, url)
					if err != nil && !errors.Is(err, context.DeadlineExceeded) {
						fmt.Println("Some error has occured: %s", err)
						return
					}
					fmt.Println(page.GetName(), "got", len(result), "proxies")
					for _, proxy := range result {
						proxy := proxy
						select {
						case chProxy <- proxy:
							continue
						case <-ctx.Done():
							break
						}
					}
				}
				fmt.Println(page.GetName(), "end")
			}()
		}
		pwg.Wait()
		close(chProxy)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		wgp := sync.WaitGroup{}
		for i := 0; i < maxConn; i++ {

			wgp.Add(1)
			go func() {
				defer wgp.Done()
				for {
					select {
					case proxy, done := <-chProxy:
						if !done {
							return
						}
						err := pinger.CheckProxy(ctx, proxy)
						if err == nil {
							pMutex.Lock()
							fmt.Println(proxies, "Working:", proxy)
							proxies += 1
							chProved <- proxy
							pMutex.Unlock()
						} else {
							//fmt.Println("Dead:", proxy, err)
						}
						continue
					case <-ctx.Done():
						return
					}
				}
			}()
		}
		wgp.Wait()
		close(chProved)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case proxy, done := <-chProved:
				if !done {
					return
				}
				err := saver.SavePage(ctx, proxy, 0)
				if err != nil {
					fmt.Println("saver error:", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()
}
