package deploy

import (
	"os/user"
	"testing"

	"github.com/teodor-pripoae/vessel/client/ssh"

	. "github.com/teodor-pripoae/vessel/Godeps/_workspace/src/gopkg.in/check.v1"
)

func TestCL(t *testing.T) { TestingT(t) }

type DeploySuite struct {
}

var _ = Suite(&DeploySuite{})

// ssh.GetConfig
func (s *DeploySuite) TestsshGetConfigNoPort(c *C) {
	server := "ssh://myuser@my.server.host"

	config, err := ssh.GetConfig(server)

	expectedConfig := ssh.Config{
		User:   "myuser",
		Server: "my.server.host",
		Key:    "/.ssh/id_rsa",
		Port:   "22",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestsshGetConfigWithPort(c *C) {
	server := "ssh://myuserb@my.server.host:2222"

	config, err := ssh.GetConfig(server)

	expectedConfig := ssh.Config{
		User:   "myuserb",
		Server: "my.server.host",
		Key:    "/.ssh/id_rsa",
		Port:   "2222",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestsshGetConfigNoUserNoPort(c *C) {
	server := "ssh://my.server.host"
	currentUser, _ := user.Current()

	config, err := ssh.GetConfig(server)

	expectedConfig := ssh.Config{
		User:   currentUser.Username,
		Server: "my.server.host",
		Port:   "22",
		Key:    "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestsshGetConfigNoUserWithPort(c *C) {
	server := "ssh://my.server.host2:2244"
	currentUser, _ := user.Current()

	config, err := ssh.GetConfig(server)

	expectedConfig := ssh.Config{
		User:   currentUser.Username,
		Server: "my.server.host2",
		Port:   "2244",
		Key:    "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestsshGetConfigNoPortService(c *C) {
	server := "ssh://myuser@my.server.host/myapp"

	config, err := ssh.GetConfig(server)

	expectedConfig := ssh.Config{
		User:    "myuser",
		Server:  "my.server.host",
		Port:    "22",
		Service: "myapp",
		Key:     "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestsshGetConfigWithPortService(c *C) {
	server := "ssh://myuserb@my.server.host:2222/myapp"

	config, err := ssh.GetConfig(server)

	expectedConfig := ssh.Config{
		User:    "myuserb",
		Server:  "my.server.host",
		Port:    "2222",
		Service: "myapp",
		Key:     "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestsshGetConfigNoUserNoPortService(c *C) {
	server := "ssh://my.server.host/foo"
	currentUser, _ := user.Current()

	config, err := ssh.GetConfig(server)

	expectedConfig := ssh.Config{
		User:    currentUser.Username,
		Server:  "my.server.host",
		Port:    "22",
		Service: "foo",
		Key:     "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

func (s *DeploySuite) TestsshGetConfigNoUserWithPortService(c *C) {
	server := "ssh://my.server.host2:2244/bar"
	currentUser, _ := user.Current()

	config, err := ssh.GetConfig(server)

	expectedConfig := ssh.Config{
		User:    currentUser.Username,
		Server:  "my.server.host2",
		Port:    "2244",
		Service: "bar",
		Key:     "/.ssh/id_rsa",
	}

	c.Assert(err, IsNil)
	c.Assert(*config, DeepEquals, expectedConfig)
}

// /ssh.GetConfig
