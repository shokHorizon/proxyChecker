package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/shokHorizon/proxyChecker/config"
	"github.com/shokHorizon/proxyChecker/internal/models"
	"github.com/shokHorizon/proxyChecker/internal/parser"
	"github.com/shokHorizon/proxyChecker/internal/pinger"
	"github.com/shokHorizon/proxyChecker/internal/provider"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	chProxy := make(chan *models.Proxy, 100)
	chProved := make(chan *models.Proxy, 100)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	ping := pinger.NewPinger(cfg.Pinger, chProxy, chProved)
	prov := provider.NewProvider(cfg.Provider)

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
						case chProxy <- &proxy:
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
		ping.Run(ctx)
		close(chProved)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		prov.Run(ctx, chProved)
	}()

	// Done ctx on wg.Wait()
	go func() {
		wg.Wait()
		cancel()
	}()

	select {
	case <-interrupt:
		fmt.Println("Interrupt signal received, exiting...")
		cancel()
	case <-ctx.Done():
		break
	}
}
