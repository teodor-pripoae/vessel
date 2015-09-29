package deploy

import (
	"fmt"

	cfg "github.com/teodor-pripoae/vessel/client/config"
	"github.com/teodor-pripoae/vessel/client/deploy/uploaders"
	"github.com/teodor-pripoae/vessel/client/ssh"
)

// Deploy is called after slug build finished
func Deploy(slugPath string, config cfg.Config, app cfg.AppConfig) error {
	for _, server := range *app.Deploy.UploadServers {
		if err := copyDeploySlug(slugPath, server, app); err != nil {
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

func copyDeploySlug(slugPath string, server string, app cfg.AppConfig) error {
	uploader, err := uploaders.NewUploader(server, app.Deploy.SlugLocation)

	if err != nil {
		return err
	}

	fmt.Printf("Uploading slug to server %v\n", uploader.Server())

	if err := uploader.Put(slugPath); err != nil {
		return err
	}

	return nil
}

func restartService(server string, app cfg.AppConfig) error {
	sshC, err := ssh.GetConfig(server)

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
