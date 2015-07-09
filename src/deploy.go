package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// Deploy is called after slug build finished
func Deploy(slugPath string, config Config, app AppConfig) error {
	for _, server := range *app.Deploy.UploadServers {
		err := copyDeploySlug(slugPath, server, app)

		if err != nil {
			return err
		}
	}

	for _, entry := range *app.Deploy.Services {
		split := strings.Split(entry, "/")

		if len(split) < 2 {
			log.Fatalf("Please enter service in the format: myuser@myapp.com/my_service")
		}
		restartService(split[0], split[1], config, app)
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

func restartService(server string, service string, config Config, app AppConfig) {
	restartCmd := fmt.Sprintf("sudo service %s restart", service)
	cmd := exec.Command("ssh", server, restartCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// returns ssh connection config
func getSSHConfig(serverURL string) (*SSHConfig, error) {
	uri, err := url.Parse(serverURL)

	if err != nil {
		log.Fatalf("Failed to parse server url %v, err: %v", serverURL, err)
		return nil, err
	}

	user, err := getSSHUser(uri)
	if err != nil {
		return nil, err
	}

	host, port, err := getSSHHostPort(uri)

	if err != nil {
		return nil, err
	}

	config := SSHConfig{
		User:   *user,
		Server: *host,
		Port:   *port,
	}

	return &config, nil
}

func getSSHUser(uri *url.URL) (*string, error) {
	if uri.User != nil {
		usr := uri.User.Username()
		return &usr, nil
	}

	usr, err := user.Current()

	if err != nil {
		return nil, err
	}

	return &usr.Username, nil
}

func getSSHHostPort(uri *url.URL) (*string, *string, error) {
	parsedHost := strings.Split(uri.Host, ":")

	if len(parsedHost) == 0 {
		log.Fatalf("server should not blank")
		return nil, nil, fmt.Errorf("Server <%v> was not valid", uri.Host)
	}

	server := parsedHost[0]
	port := "22"

	if len(parsedHost) >= 2 {
		port = parsedHost[1]
	}

	return &server, &port, nil
}
