package main

import (
	"log"

	"github.com/go-resty/resty/v2"
)

type hcResult struct {
	IsHealthy bool
}

func (s *Site) check(prod bool) {
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