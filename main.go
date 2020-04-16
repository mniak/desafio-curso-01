package main

import (
	"log"
	"sync"

	"github.com/go-resty/resty/v2"
)

type hcResult struct {
	IsHealthy bool
}

func check(wg *sync.WaitGroup, s Site, prod bool) {
	defer wg.Done()

	client := resty.New()
	resp, err := client.R().
		SetResult(&hcResult{}).
		Get(s.url(prod))

	if err != nil {
		log.Printf("Erro ao testar site  %v: %v\n", s.Name, err)
		return
	}

	result := resp.Result().(*hcResult)
	if result.IsHealthy {
		log.Printf("%v está saudável\n", s.Name)
	} else {
		log.Printf("%v está com erro\n", s.Name)
	}
}

func main() {

	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	sites := config.Sites
	prod := isProd()

	var wg sync.WaitGroup
	for _, s := range sites {
		wg.Add(1)
		go check(&wg, s, prod)
	}
	wg.Wait()
}
