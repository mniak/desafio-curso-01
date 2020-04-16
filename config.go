package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

//Site represents one site for which run the health check
type Site struct {
	Name      string
	Domain    string
	SubDomain string
}

func (s *Site) url(isprod bool) string {
	if isprod {
		return fmt.Sprintf("https://%s.%s/healthcheck", s.SubDomain, s.Domain)
	}
	return fmt.Sprintf("https://%ssandbox.%s/healthcheck", s.SubDomain, s.Domain)
}

//Config represents the configuration file
type Config struct {
	Sites []Site
}

func readConfig() (config Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if viper.Unmarshal(&config); err != nil {
		return
	}

	return
}

func isProd() bool {

	prodPtr := flag.Bool("prod", false, "Use production environment")
	sboxPtr := flag.Bool("sbox", false, "Use sandbox environment")
	flag.Parse()

	flagProd, flagSbox := *prodPtr, *sboxPtr

	if flagSbox || flagProd {
		return flagProd && !flagSbox
	}
	prompt := promptui.Select{
		Label:     "Ambiente",
		Items:     []string{"Sandbox", "Produção"},
		IsVimMode: false,
	}
	n, _, err := prompt.Run()
	if err != nil {
		log.Printf("Erro na leitura %v. Usando sandbox\n", err)
		n = 0
	}
	return n == 1
}
