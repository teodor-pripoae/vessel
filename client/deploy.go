package client

import "fmt"

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

	sshC, err := getSSHConfig(server)

	if err != nil {
		return err
	}

	fmt.Printf("Uploading slug to server %v\n", sshC.Server)

	if err := sshC.Scp(slugPath, *app.Deploy.SlugLocation); err != nil {
		return err
	}

	return nil
}

func restartService(server string, app AppConfig) error {
	sshC, err := getSSHConfig(server)

	if err != nil {
		return err
	}

	fmt.Printf("Restarting service %v on server %v\n", sshC.Service, sshC.Server)

	restartCmd := fmt.Sprintf("sudo service %s restart", sshC.Service)

	stdout, stderr, err := sshC.Run(restartCmd)

	if err != nil {
		fmt.Printf("Stdout: %s\n", stdout)
		fmt.Printf("Stderr: %s\n", stderr)
		return err
	}

	return nil
}
