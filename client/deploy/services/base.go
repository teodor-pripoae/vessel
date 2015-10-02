package services

import (
	"errors"
	"strings"
)

const (
	// SSHServiceType is returned by SSHService.Type()
	SSHServiceType = 1
	// MarathonServiceType is returned by MarathonService.Type()
	MarathonServiceType = 2
)

// Service - base interface for defining a service manager
type Service interface {
	Restart() error
	Type() int
	Service() string
	Server() string
}

var (
	// ErrServiceNotFound Only ssh/upstart and marathon apps supported
	ErrServiceNotFound = errors.New("Cannot find service protocol, ssh:// and mrt:// supported")
)

// NewService returns new service implementation
func NewService(server string) (Service, error) {
	if strings.HasPrefix(server, "ssh://") {
		return NewSSHService(server)
	}

	if strings.HasPrefix(server, "mrt://") {
		return NewMarathonService(server)
	}

	return nil, ErrServiceNotFound
}
