package services

import (
	"fmt"

	"github.com/teodor-pripoae/vessel/client/ssh"
)

// SSHService defines a service managed using ssh
type SSHService struct {
	sshC    ssh.Config
	service string
	server  string
}

// NewSSHService returns new ssh service
func NewSSHService(server string) (*SSHService, error) {
	sshC, err := ssh.GetConfig(server)

	if err != nil {
		return nil, err
	}

	service := SSHService{
		service: sshC.Service,
		server:  sshC.Server,
		sshC:    *sshC,
	}

	return &service, nil
}

// Restart service using SSH
func (c *SSHService) Restart() error {
	fmt.Printf("Restarting service %s on server %s\n", c.service, c.server)
	restartCmd := fmt.Sprintf("sudo service %s restart", c.service)

	stdout, stderr, err := c.sshC.Run(restartCmd)

	if err != nil {
		fmt.Printf("Stdout: %s\n", stdout)
		fmt.Printf("Stderr: %s\n", stderr)
		return err
	}

	return nil
}

// Type returns service type
func (c *SSHService) Type() int {
	return SSHServiceType
}

// Service returns service name
func (c *SSHService) Service() string {
	return c.service
}

// Server returns server name
func (c *SSHService) Server() string {
	return c.server
}
