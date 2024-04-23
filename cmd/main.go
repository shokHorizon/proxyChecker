package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/shokHorizon/proxyChecker/config"
	"github.com/shokHorizon/proxyChecker/internal/models"
	"github.com/shokHorizon/proxyChecker/internal/parser"
	"github.com/shokHorizon/proxyChecker/internal/pinger"
	"sync"
	"time"
)

const (
	maxConn    = 5000
	maxWorkers = 500
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	chProxy := make(chan models.Proxy, 400)
	//chProved := make(chan string, 400)
	pMutex := sync.Mutex{}
	proxies := 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		pwg := sync.WaitGroup{}
		for _, sourceCfg := range cfg.Sources {
			source := parser.NewFromConfig(sourceCfg)
			pwg.Add(1)
			go func() {
				defer pwg.Done()
				for _, url := range source.GetUrls() {
					ctx, _ := context.WithTimeout(ctx, time.Second*10)
					result, err := source.Parse(ctx, url)
					if err != nil && !errors.Is(err, context.DeadlineExceeded) {
						fmt.Println("Some error has occured: %s", err)
						return
					}
					fmt.Println(source.GetName(), "got", len(result), "proxies")
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
				fmt.Println(source.GetName(), "end")
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
						err := pinger.CheckProxy(ctx, &proxy)
						if err != nil {
							err = pinger.CheckProxySocks(ctx, &proxy)
							if err != nil {
								err = pinger.CheckProxySocks(ctx, &proxy)
							}
						}
						if err == nil {
							pMutex.Lock()
							fmt.Println(proxies, proxy.Origin, "Working:", proxy.Url())
							proxies += 1
							//chProved <- provider
							pMutex.Unlock()
						} else {
							//fmt.Println("Dead:", provider, err)
						}
						continue
					case <-ctx.Done():
						return
					}
				}
			}()
		}
		wgp.Wait()
		//close(chProved)
	}()

	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	for {
	//		select {
	//		case provider, done := <-chProved:
	//			if !done {
	//				return
	//			}
	//			err := saver.SavePage(ctx, provider, 0)
	//			if err != nil {
	//				fmt.Println("saver error:", err)
	//			}
	//		case <-ctx.Done():
	//			return
	//		}
	//	}
	//}()

	wg.Wait()
}
