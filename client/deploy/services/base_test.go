package services

import (
	"testing"

	. "github.com/teodor-pripoae/vessel/Godeps/_workspace/src/gopkg.in/check.v1"
)

func TestServices(t *testing.T) { TestingT(t) }

type ServicesSuite struct {
}

var _ = Suite(&ServicesSuite{})

// Test if using ssh url it returns a ssh service
func (s *ServicesSuite) TestNewServiceWithSshURL(c *C) {
	server := "ssh://myuser@my.server.host/myapp"

	service, err := NewService(server)

	c.Assert(err, IsNil)
	c.Assert(service, NotNil)

	c.Assert(service.Type(), Equals, SSHServiceType)
	c.Assert(service.Server(), Equals, "my.server.host")
	c.Assert(service.Service(), Equals, "myapp")
}

// Test if using marathon url returns correct values
func (s *ServicesSuite) TestNewServiceWithMarathonURL(c *C) {
	server := "mrt://my.server.com:8080/my-app"

	service, err := NewService(server)

	c.Assert(err, IsNil)
	c.Assert(service.Type(), Equals, MarathonServiceType)

	c.Assert(service.Server(), Equals, "http://my.server.com:8080")
	c.Assert(service.Service(), Equals, "my-app")
}

// Test if using unsupported url returns error
func (s *ServicesSuite) TestNewServiceWithUnsupportedURL(c *C) {
	server := "http://myuser@my.server.host"

	service, err := NewService(server)

	c.Assert(err, Equals, ErrServiceNotFound)
	c.Assert(service, IsNil)
}
