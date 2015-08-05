package client

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
		Key:    "/.ssh/id_rsa",
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
		Key:    "/.ssh/id_rsa",
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
		Key:    "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestGetSSHConfigNoUserWithPort(c *C) {
	server := "ssh://my.server.host2:2244"
	currentUser, _ := user.Current()

	config, err := getSSHConfig(server)

	expectedConfig := SSHConfig{
		User:   currentUser.Username,
		Server: "my.server.host2",
		Port:   "2244",
		Key:    "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestGetSSHConfigNoPortService(c *C) {
	server := "ssh://myuser@my.server.host/myapp"

	config, err := getSSHConfig(server)

	expectedConfig := SSHConfig{
		User:    "myuser",
		Server:  "my.server.host",
		Port:    "22",
		Service: "myapp",
		Key:     "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestGetSSHConfigWithPortService(c *C) {
	server := "ssh://myuserb@my.server.host:2222/myapp"

	config, err := getSSHConfig(server)

	expectedConfig := SSHConfig{
		User:    "myuserb",
		Server:  "my.server.host",
		Port:    "2222",
		Service: "myapp",
		Key:     "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestGetSSHConfigNoUserNoPortService(c *C) {
	server := "ssh://my.server.host/foo"
	currentUser, _ := user.Current()

	config, err := getSSHConfig(server)

	expectedConfig := SSHConfig{
		User:    currentUser.Username,
		Server:  "my.server.host",
		Port:    "22",
		Service: "foo",
		Key:     "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestGetSSHConfigNoUserWithPortService(c *C) {
	server := "ssh://my.server.host2:2244/bar"
	currentUser, _ := user.Current()

	config, err := getSSHConfig(server)

	expectedConfig := SSHConfig{
		User:    currentUser.Username,
		Server:  "my.server.host2",
		Port:    "2244",
		Service: "bar",
		Key:     "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

// /getSSHConfig
