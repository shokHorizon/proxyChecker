package saver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const steamUrl = `https://steamcommunity.com/market/search/render/?query=&start=%d&count=100&search_descriptions=0&sort_column=popular&sort_dir=desc&appid=730&norender=1`

func SavePage(ctx context.Context, proxyUrl string, index int) error {
	proxy, err := url.Parse("http://" + proxyUrl)
	if err != nil {
		return err
	}

	transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
	client := &http.Client{Transport: transport}

	ctx, _ = context.WithTimeout(ctx, time.Second)

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(steamUrl, index+1), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// bytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	// content := string(bytes)

	// re := regexp.MustCompile(`"Смартфон Apple iPhone 15 128 Гб nano-SIM \+ eSIM Blue"`)
	// matches := re.FindAllString(content, -1)

	// if len(matches) > 0 {
	// 	fmt.Printf("%s got %d results\n", proxyUrl, len(matches))
	// }

	// Create file
	f, err := os.Create(fmt.Sprintf("page%d.html", index))
	if err != nil {
		return err
	}
	defer f.Close()

	// Write body to file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("page", index, "got", resp.StatusCode)

	return nil
}
