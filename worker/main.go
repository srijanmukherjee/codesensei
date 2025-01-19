package main

import (
	"log"

	"github.com/srijanmukherjee/codesensei/shared/config"
)

func main() {
	config := config.LoadConfigFile("config/codesensei.dev.yaml")
	log.Printf("%+v", config)
}
