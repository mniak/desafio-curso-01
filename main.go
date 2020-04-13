package main

import (
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/manifoldco/promptui"
)

type hcResult struct {
	IsHealthy bool
}

type site struct {
	Name      string
	Domain    string
	SubDomain string
}

func (s *site) url(sandbox bool) string {
	if sandbox {
		return fmt.Sprintf("https://%ssandbox.%s/healthcheck", s.SubDomain, s.Domain)
	}
	return fmt.Sprintf("https://%s.%s/healthcheck", s.SubDomain, s.Domain)
}

func check(wg *sync.WaitGroup, s site, sandbox bool) {
	defer wg.Done()

	client := resty.New()
	resp, err := client.R().
		SetResult(&hcResult{}).
		Get(s.url(sandbox))

	if err != nil {
		fmt.Printf("Erro ao testar site  %v: %v\n", s.Name, err)
		return
	}

	result := resp.Result().(*hcResult)
	if result.IsHealthy {
		fmt.Printf("%v está saudável\n", s.Name)
	} else {
		fmt.Printf("%v está com erro\n", s.Name)
	}
}

func askSandbox() bool {

	prompt := promptui.Select{
		Label: "Ambiente",
		Items: []string{"Sandbox", "Produção"},
	}
	n, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Erro na leitura %v. Usando sandbox\n", err)
		n = 1
	}
	fmt.Printf("N selecionado: %d\n", n)
	sandbox := n == 0
	return sandbox
}

func main() {

	sites := []site{
		{
			Name:      "Omni API",
			Domain:    "cieloecommerce.cielo.com.br",
			SubDomain: "omni",
		},
		{
			Name:      "Omni Query API",
			Domain:    "cieloecommerce.cielo.com.br",
			SubDomain: "omniquery",
		},
		{
			Name:      "Parameters Download API",
			Domain:    "cieloecommerce.cielo.com.br",
			SubDomain: "parametersdownload",
		},
		{
			Name:      "Merchants API",
			Domain:    "cieloecommerce.cielo.com.br",
			SubDomain: "merchantapi",
		},
		{
			Name:      "Onboarding API",
			Domain:    "cieloecommerce.cielo.com.br",
			SubDomain: "onboarding",
		},
	}

	sandbox := askSandbox()

	var wg sync.WaitGroup
	for _, s := range sites {
		wg.Add(1)
		go check(&wg, s, sandbox)
	}

	wg.Wait()
	fmt.Println("-- That's all folks! --")
}
