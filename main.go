package main

import (
	"log"
	"sync"
)

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
		go func(s Site, wg *sync.WaitGroup) {
			defer wg.Done()
			s.check(prod)
		}(s, &wg)
	}
	wg.Wait()
}
