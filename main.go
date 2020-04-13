package main

import (
	"log"
	"sync"

	"github.com/go-resty/resty/v2"
)

type hcResult struct {
	IsHealthy bool
}

func check(wg *sync.WaitGroup, url string) {
	defer wg.Done()

	client := resty.New()
	resp, err := client.R().
		SetResult(&hcResult{}).
		Get(url)

	// resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error while fetching hc %v: %v\n", url, err)
		return
	}

	result := resp.Result().(*hcResult)
	if result.IsHealthy {
		log.Printf("Site %v is healthy\n", url)
	} else {
		log.Printf("Site %v is not healthy\n", url)
	}
}

func main() {

	urls := []string{
		"https://parametersdownloadsandbox.cieloecommerce.cielo.com.br/healthcheck",
		"https://merchantapisandbox.cieloecommerce.cielo.com.br/healthcheck",
		"https://omnisandbox.cieloecommerce.cielo.com.br/healthcheck",
		"https://omniquerysandbox.cieloecommerce.cielo.com.br/healthcheck",
		"https://onboardingsandbox.cieloecommerce.cielo.com.br/healthcheck",
	}
	var wg sync.WaitGroup

	for _, d := range urls {
		wg.Add(1)
		go check(&wg, d)
	}

	wg.Wait()
}
