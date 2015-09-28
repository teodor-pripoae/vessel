package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/teodor-pripoae/vessel/Godeps/_workspace/src/github.com/BurntSushi/toml"
	"github.com/teodor-pripoae/vessel/client/build"
	cfg "github.com/teodor-pripoae/vessel/client/config"
	"github.com/teodor-pripoae/vessel/client/deploy"
	"github.com/teodor-pripoae/vessel/client/notify"
)

func getConfig(args []string) cfg.Config {
	if len(args) < 4 {
		log.Fatalf(fmt.Sprintf("Required at least 3 parameters, app, commit, author, only %v provided", len(args)))
	}

	return cfg.Config{Config: args[1], Commit: args[2], Deployer: args[3]}
}

func getAppConfigData(location string) string {
	dat, err := ioutil.ReadFile(location)

	if err != nil {
		log.Fatalf("Failed to read from file %v", location)
	}

	return string(dat)
}

func getAppConfig(c cfg.Config) cfg.AppConfig {
	var appConfig cfg.AppConfig
	if _, err := toml.Decode(getAppConfigData(c.Config), &appConfig); err != nil {
		log.Fatalf("Error decoding config %v", err)
	}
	return appConfig
}

func main() {
	config := getConfig(os.Args)
	appConfig := getAppConfig(config)

	notify.OnStart(config, appConfig)

	slugPath := build.Build(config, appConfig)
	err := deploy.Deploy(slugPath, config, appConfig)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		notify.OnFailure(config, appConfig)
	} else {
		notify.OnSuccess(config, appConfig)
	}
}
