package main

import (
	"os/user"
	"testing"

	. "github.com/kuende/api_gateway/Godeps/_workspace/src/gopkg.in/check.v1"
)

func TestCL(t *testing.T) { TestingT(t) }

type DeploySuite struct {
}

var _ = Suite(&DeploySuite{})

// getSSHConfig
func (s *DeploySuite) TestGetSSHConfigNoPort(c *C) {
	server := "ssh://myuser@my.server.host"

	config, err := getSSHConfig(server)

	expectedConfig := SSHConfig{
		User:   "myuser",
		Server: "my.server.host",
		Port:   "22",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestGetSSHConfigWithPort(c *C) {
	server := "ssh://myuserb@my.server.host:2222"

	config, err := getSSHConfig(server)

	expectedConfig := SSHConfig{
		User:   "myuserb",
		Server: "my.server.host",
		Port:   "2222",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestGetSSHConfigNoUserNoPort(c *C) {
	server := "ssh://my.server.host"
	currentUser, _ := user.Current()

	config, err := getSSHConfig(server)

	expectedConfig := SSHConfig{
		User:   currentUser.Username,
		Server: "my.server.host",
		Port:   "22",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

// /getSSHConfig
