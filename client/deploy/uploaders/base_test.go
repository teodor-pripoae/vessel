package uploaders

import (
	"testing"

	. "github.com/teodor-pripoae/vessel/Godeps/_workspace/src/gopkg.in/check.v1"
)

func TestUploaders(t *testing.T) { TestingT(t) }

type UploadersSuite struct {
}

var _ = Suite(&UploadersSuite{})

// Test if using ssh url and upload location returns ssh uploader
func (s *UploadersSuite) TestNewUploaderWithSSHURL(c *C) {
	server := "ssh://myuser@my.server.host"
	uploadLocation := "/etc/vessel/deploy/slugs/myapp/current.tgz"

	uploader, err := NewUploader(server, &uploadLocation)

	c.Assert(err, IsNil)
	c.Assert(uploader, NotNil)

	c.Assert(uploader.Type(), Equals, SSHUploaderType)
	c.Assert(uploader.remoteLocation(), Equals, uploadLocation)
	c.Assert(uploader.Server(), Equals, "my.server.host")
}

// Test if using ssh url and no upload location returns error
func (s *UploadersSuite) TestNewUploaderWithSSHURLAndNoLocation(c *C) {
	server := "ssh://myuser@my.server.host"

	uploader, err := NewUploader(server, nil)

	c.Assert(err, Equals, ErrLocationNotSet)
	c.Assert(uploader, IsNil)
}

// Test if using webdav url returns correct values
func (s *UploadersSuite) TestNewUploaderWithWebDavURL(c *C) {
	server := "webdav://my.server.com:8000/slugs/myapp/current.tgz"
	uploadLocation := "/etc/vessel/deploy/slugs/myapp/current.tgz"

	uploader, err := NewUploader(server, &uploadLocation)

	c.Assert(err, IsNil)
	c.Assert(uploader.remoteLocation(), Equals, "slugs/myapp/current.tgz")
	c.Assert(uploader.Server(), Equals, "http://my.server.com:8000/slugs/myapp/current.tgz")
}

// Test if using unsupported url returns error
func (s *UploadersSuite) TestNewUploaderWithUnsupportedURL(c *C) {
	server := "http://myuser@my.server.host"
	uploadLocation := "/etc/vessel/deploy/slugs/myapp/current.tgz"

	uploader, err := NewUploader(server, &uploadLocation)

	c.Assert(err, Equals, ErrUploaderNotFound)
	c.Assert(uploader, IsNil)
}
