package main

// taken from github.com/hypersleep/easyssh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"os/user"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHConfig keeps ssh session stuff
type SSHConfig struct {
	User     string
	Server   string
	Port     string
	Key      string
	Password string
	Service  string
}

func getKeyFile(keyPath string) (ssh.Signer, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	// TODO: use strings.Join
	file := usr.HomeDir + keyPath
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	pubkey, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

// connects to remote server using MakeConfig struct and returns *ssh.Session
func (ssh_conf *SSHConfig) connect() (*ssh.Session, error) {
	// auths holds the detected ssh auth methods
	auths := []ssh.AuthMethod{}

	// figure out what auths are requested, what is supported
	if ssh_conf.Password != "" {
		auths = append(auths, ssh.Password(ssh_conf.Password))
	}

	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers))
		defer sshAgent.Close()
	}

	if pubkey, err := getKeyFile(ssh_conf.Key); err == nil {
		auths = append(auths, ssh.PublicKeys(pubkey))
	}

	config := &ssh.ClientConfig{
		User: ssh_conf.User,
		Auth: auths,
	}

	client, err := ssh.Dial("tcp", ssh_conf.Server+":"+ssh_conf.Port, config)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Run executes command on remote machine and returns STDOUT
func (ssh_conf *SSHConfig) Run(command string) (string, error) {
	session, err := ssh_conf.connect()

	if err != nil {
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run(command)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// Scp uploads sourceFile to remote machine like native scp console app.
func (ssh_conf *SSHConfig) Scp(sourceFile string, destFile string) error {
	session, err := ssh_conf.connect()

	if err != nil {
		return err
	}
	defer session.Close()

	src, srcErr := os.Open(sourceFile)

	if srcErr != nil {
		return srcErr
	}

	srcStat, statErr := src.Stat()

	if statErr != nil {
		return statErr
	}

	go func() {
		w, _ := session.StdinPipe()

		fmt.Fprintln(w, "C0644", srcStat.Size(), destFile)

		if srcStat.Size() > 0 {
			io.Copy(w, src)
			fmt.Fprint(w, "\x00")
			w.Close()
		} else {
			fmt.Fprint(w, "\x00")
			w.Close()
		}
	}()

	if err := session.Run(fmt.Sprintf("scp -t %s", destFile)); err != nil {
		return err
	}

	return nil
}

// returns ssh connection config
func getSSHConfig(serverURL string) (*SSHConfig, error) {
	uri, err := url.Parse(serverURL)

	if err != nil {
		log.Fatalf("Failed to parse server url %v, err: %v", serverURL, err)
		return nil, err
	}

	user, err := getSSHUser(uri)
	if err != nil {
		return nil, err
	}

	host, port, err := getSSHHostPort(uri)

	if err != nil {
		return nil, err
	}

	service, err := getSSHService(uri)

	if err != nil {
		return nil, err
	}

	config := SSHConfig{
		User:    *user,
		Server:  *host,
		Port:    *port,
		Service: service,
	}

	return &config, nil
}

func getSSHUser(uri *url.URL) (*string, error) {
	if uri.User != nil {
		usr := uri.User.Username()
		return &usr, nil
	}

	usr, err := user.Current()

	if err != nil {
		return nil, err
	}

	return &usr.Username, nil
}

func getSSHHostPort(uri *url.URL) (*string, *string, error) {
	parsedHost := strings.Split(uri.Host, ":")

	if len(parsedHost) == 0 {
		log.Fatalf("server should not blank")
		return nil, nil, fmt.Errorf("Server <%v> was not valid", uri.Host)
	}

	server := parsedHost[0]
	port := "22"

	if len(parsedHost) >= 2 {
		port = parsedHost[1]
	}

	return &server, &port, nil
}

func getSSHService(uri *url.URL) (string, error) {
	path := strings.Trim(uri.Path, "/")
	splits := strings.Split(path, "/")

	if len(splits) == 0 {
		return "", nil
	}

	if len(splits) > 1 {
		return "", fmt.Errorf("Service not well formatted: %v", uri.Path)
	}

	return splits[0], nil
}
