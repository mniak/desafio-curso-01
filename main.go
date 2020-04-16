package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

type hcResult struct {
	IsHealthy bool
}

func (s *site) url(prod bool) string {
	if prod {
		return fmt.Sprintf("https://%s.%s/healthcheck", s.SubDomain, s.Domain)
	}
	return fmt.Sprintf("https://%ssandbox.%s/healthcheck", s.SubDomain, s.Domain)
}

func check(wg *sync.WaitGroup, s site, prod bool) {
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

func determineIfProd(flagSbox, flagProd bool) bool {

	if flagSbox || flagProd {
		return flagProd && !flagSbox
	}
	prompt := promptui.Select{
		Label:     "Ambiente",
		Items:     []string{"Sandbox", "Produção"},
		IsVimMode: false,
	}
	n, _, err := prompt.Run()
	log.Println(n)
	if err != nil {
		log.Printf("Erro na leitura %v. Usando sandbox\n", err)
		n = 0
	}
	return n == 1
}

type site struct {
	Name      string
	Domain    string
	SubDomain string
}

//Config represents the configuration file
type Config struct {
	Sites []struct {
		Name      string
		Domain    string
		SubDomain string
	}
}

func main() {

	var err error

	prodPtr := flag.Bool("prod", false, "Use production environment")
	sboxPtr := flag.Bool("sbox", false, "Use sandbox environment")
	flag.Parse()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err = viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	var config Config
	if viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	sites := config.Sites

	prod := determineIfProd(*sboxPtr, *prodPtr)

	var wg sync.WaitGroup
	for _, s := range sites {
		wg.Add(1)
		go check(&wg, s, prod)
	}

	wg.Wait()
	log.Println("-- That's all folks! --")
}
