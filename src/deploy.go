package main

import (
	"fmt"
	"os"
)

// Deploy is called after slug build finished
func Deploy(slugPath string, config Config, app AppConfig) error {
	for _, server := range *app.Deploy.UploadServers {
		err := copyDeploySlug(slugPath, server, app)

		if err != nil {
			return err
		}
	}

	for _, service := range *app.Deploy.Services {
		if err := restartService(service, app); err != nil {
			return err
		}
	}

	return nil
}

func copyDeploySlug(slugPath string, server string, app AppConfig) error {
	if app.Deploy.SlugLocation == nil {
		return fmt.Errorf("Deploy slug location not defined")
	}

	sshConfig, err := getSSHConfig(server)

	if err != nil {
		return err
	}

	if err := sshConfig.Scp(slugPath, *app.Deploy.SlugLocation); err != nil {
		return err
	}

	return nil
}

func restartService(server string, app AppConfig) error {
	sshC, err := getSSHConfig(server)

	if err != nil {
		return err
	}

	restartCmd := fmt.Sprintf("sudo service %s restart", sshC.Service)

	output, err := sshC.Run(restartCmd)

	fmt.Fprintf(os.Stderr, output)

	if err != nil {
		return err
	}

	return nil
}
