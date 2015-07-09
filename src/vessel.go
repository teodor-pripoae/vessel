package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func getConfig(args []string) Config {
	if len(args) < 4 {
		log.Fatalf(fmt.Sprintf("Required at least 3 parameters, app, commit, author, only %v provided", len(args)))
	}

	config := Config{Config: args[1], Commit: args[2], Deployer: args[3]}
	return config
}

func getAppConfigData(location string) string {
	dat, err := ioutil.ReadFile(location)

	if err != nil {
		log.Fatalf("Failed to read from file %v", location)
	}

	return string(dat)
}

func getAppConfig(c Config) AppConfig {
	var appConfig AppConfig
	if _, err := toml.Decode(getAppConfigData(c.Config), &appConfig); err != nil {
		log.Fatalf("Error decoding config %v", err)
	}
	return appConfig
}

func main() {
	config := getConfig(os.Args)
	appConfig := getAppConfig(config)

	NotifyOnStart(config, appConfig)

	slugPath := Build(config, appConfig)
	err := Deploy(slugPath, config, appConfig)

	if err != nil {
		fmt.Printf("Error: %v", err)
		NotifyOnFailure(config, appConfig)
	} else {
		NotifyOnSuccess(config, appConfig)
	}
}
