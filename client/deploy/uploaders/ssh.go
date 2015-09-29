package uploaders

import "github.com/teodor-pripoae/vessel/client/ssh"

// SSHUploader uploads a local path with ssh
type SSHUploader struct {
	SSHConfig      ssh.Config
	RemoteLocation string
}

// NewSSHUploader returns new webdav uploader
func NewSSHUploader(server string, uploadLocation *string) (Uploader, error) {
	if uploadLocation == nil {
		return nil, ErrLocationNotSet
	}

	sshC, err := ssh.GetConfig(server)

	if err != nil {
		return nil, err
	}

	uploader := SSHUploader{
		SSHConfig:      *sshC,
		RemoteLocation: *uploadLocation,
	}

	return &uploader, nil
}

// Put localPath to remote location using ssh config
func (u *SSHUploader) Put(localPath string) error {
	return u.SSHConfig.Scp(localPath, u.RemoteLocation)
}

// Type returns uploader type
func (u *SSHUploader) Type() int {
	return SSHUploaderType
}

// Server returns server name
func (u *SSHUploader) Server() string {
	return u.SSHConfig.Server
}

func (u *SSHUploader) remoteLocation() string {
	return u.RemoteLocation
}
