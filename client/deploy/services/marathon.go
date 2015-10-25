package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// MarathonService defines a marathon app
type MarathonService struct {
	url string
	app string
}

// NewMarathonService returns a new marathon app service
func NewMarathonService(server string) (*MarathonService, error) {
	uri, err := url.Parse(server)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s", uri.Host)
	app := strings.Trim(uri.Path, "/")

	service := MarathonService{
		url: url,
		app: app,
	}

	return &service, nil
}

// Restart a marathon app
func (c *MarathonService) Restart() error {
	fmt.Printf("Restarting app %s on server %s\n", c.app, c.url)
	hc := http.Client{}
	url := fmt.Sprintf("%s/v2/apps/%s/restart?force=true", c.url, c.app)
	req, err := http.NewRequest("POST", url, strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	_, err = hc.Do(req)

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(resp.Body)
	// fmt.Println(buf.String())

	return err
}

// Type returns service type
func (c *MarathonService) Type() int {
	return MarathonServiceType
}

// Service returns app name
func (c *MarathonService) Service() string {
	return c.app
}

// Server returns server name
func (c *MarathonService) Server() string {
	return c.url
}
