package uploaders

import (
	"errors"
	"strings"
)

const (
	// SSHUploaderType is returned by SSHUploader.Type()
	SSHUploaderType = 1
	// WebDavUploaderType is returned by WebDavUploader.Type()
	WebDavUploaderType = 2
)

var (
	//ErrLocationNotSet means slug_location was not set
	ErrLocationNotSet = errors.New("Deploy slug location not defined")
	// ErrUploaderNotFound means upload protocol is not supported
	ErrUploaderNotFound = errors.New("Cannot find uploader protocol, ssh:// and webdav:// supported")
)

// Uploader interface implemented by all uploaders
type Uploader interface {
	Put(string) error
	Type() int
	Server() string
	remoteLocation() string
}

// NewUploader returns new generic uploader
func NewUploader(server string, uploadLocation *string) (Uploader, error) {
	if strings.HasPrefix(server, "ssh://") {
		return NewSSHUploader(server, uploadLocation)
	}

	if strings.HasPrefix(server, "webdav://") {
		return NewWebDavUploader(server)
	}

	return nil, ErrUploaderNotFound
}
