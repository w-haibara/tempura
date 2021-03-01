package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"tempura/config"
	"tempura/generate"
	"tempura/serve"
)

var (
	configFile = flag.String("f", "config.json", "path to the configuration file")
	servePort  = flag.String("serve", "", "cert.pem file path")
)

func run() int {
	log.Println("--- --- ---")

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

	if *servePort != "" {
		log.Println("serve:", *servePort)
		serve.Serve(*servePort)
		return 0
	}

	return 0
}

func main() {
	flag.Parse()
	os.Exit(run())
}
