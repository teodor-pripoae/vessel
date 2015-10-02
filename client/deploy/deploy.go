package deploy

import (
	"fmt"

	cfg "github.com/teodor-pripoae/vessel/client/config"
	"github.com/teodor-pripoae/vessel/client/deploy/services"
	"github.com/teodor-pripoae/vessel/client/deploy/uploaders"
)

// Deploy is called after slug build finished
func Deploy(slugPath string, config cfg.Config, app cfg.AppConfig) error {
	for _, server := range *app.Deploy.UploadServers {
		if err := uploadSlug(slugPath, server, app); err != nil {
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

// Upload slug to servers using ssh or upload to webdav
func uploadSlug(slugPath string, server string, app cfg.AppConfig) error {
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

// Restart services remote using ssh or marathon API
func restartService(server string, app cfg.AppConfig) error {
	service, err := services.NewService(server)

	if err != nil {
		return err
	}

	if err = service.Restart(); err != nil {
		return err
	}

	return nil
}
