package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Deploy is called after slug build finished
func Deploy(slugPath string, config Config, app AppConfig) {
	for _, server := range *app.Deploy.UploadServers {
		copyDeploySlug(slugPath, server, config, app)
	}

	for _, entry := range *app.Deploy.Services {
		split := strings.Split(entry, "/")

		if len(split) < 2 {
			log.Fatalf("Please enter service in the format: myuser@myapp.com/my_service")
		}
		restartService(split[0], split[1], config, app)
	}
}

func copyDeploySlug(slugPath string, server string, config Config, app AppConfig) {
	if app.Deploy.SlugLocation == nil {
		log.Fatalf("Deploy slug location not defined")
	}
	destination := fmt.Sprintf("%s:%s", server, *app.Deploy.SlugLocation)

	cmd := exec.Command("scp", slugPath, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func restartService(server string, service string, config Config, app AppConfig) {
	restartCmd := fmt.Sprintf("sudo service %s restart", service)
	cmd := exec.Command("ssh", server, restartCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
