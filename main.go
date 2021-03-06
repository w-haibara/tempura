package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/w-haibara/tempura/config"
	"github.com/w-haibara/tempura/generate"
	"github.com/w-haibara/tempura/serve"
)

var (
	configFile = flag.String("f", "config.json", "path to the configuration file")
	servePort  = flag.String("serve", "", "cert.pem file path")
)

func run() int {
	log.Println("--- starting ---")

	log.Printf("config file: %v\n", *configFile)

	json, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Panic(err)
	}

	c, err := config.Configure(json)
	if err != nil {
		log.Panic(err)
	}

	if err := generate.Generate(c); err != nil {
		log.Panic(err)
	}

	log.Println("--- completed ---")

	if *servePort != "" {
		log.Println("--- serve:", *servePort, "---")
		serve.Serve(*servePort)
		log.Println("--- server stopped ---")
		return 0
	}

	return 0
}

func main() {
	flag.Parse()
	os.Exit(run())
}
