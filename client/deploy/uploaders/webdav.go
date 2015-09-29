package uploaders

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/teodor-pripoae/vessel/client/deploy/uploaders/dav"
)

// WebDavUploader uploads a local path to a webdav server
type WebDavUploader struct {
	Conn           dav.Session
	RemoteLocation string
	server         string
}

// NewWebDavUploader returns new webdav uploader
func NewWebDavUploader(server string) (Uploader, error) {
	uri, err := url.Parse(server)

	if err != nil {
		return nil, err
	}

	serverName := fmt.Sprintf("http://%s", uri.Host)

	s, err := dav.NewSession(serverName)

	if err != nil {
		return nil, err
	}

	location := strings.Trim(uri.Path, "/")

	uploader := WebDavUploader{
		Conn:           *s,
		RemoteLocation: location,
		server:         serverName,
	}

	return &uploader, nil
}

// Put localPath to a remote webdav folder
func (u *WebDavUploader) Put(localPath string) error {
	file, err := ioutil.ReadFile(localPath)
	if err != nil {
		return err
	}
	return u.Conn.Put(u.RemoteLocation, file)
}

// Type returns uploader type
func (u *WebDavUploader) Type() int {
	return WebDavUploaderType
}

// Server returns server name
func (u *WebDavUploader) Server() string {
	return fmt.Sprintf("%s/%s", u.server, u.RemoteLocation)
}

func (u *WebDavUploader) remoteLocation() string {
	return u.RemoteLocation
}
