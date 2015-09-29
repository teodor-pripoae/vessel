package dav

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
)

// Session keeps connection info
type Session struct {
	base     *url.URL
	client   http.Client
	username string
	password string
}

// NewSession returns new webdav session
func NewSession(rooturl string) (s *Session, err error) {
	jar, _ := cookiejar.New(nil)
	cl := http.Client{Jar: jar}
	u, _ := url.Parse(rooturl)
	s = &Session{u, cl, "", ""}
	return
}

// SetBasicAuth sets basic authentication for session
func (s *Session) SetBasicAuth(user, pass string) {
	s.username = user
	s.password = pass
	return
}

// Put a file to webdav server
func (s *Session) Put(name string, data []byte) (err error) {
	req, err := http.NewRequest("PUT", s.abs(name), bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Host = s.base.Host
	if s.username != "" {
		req.SetBasicAuth(s.username, s.password)
	}
	req.ContentLength = int64(len(data))
	res, err := s.doRequest(req)
	if err != nil {
		return err
	}
	err = s.res2Err(res, []int{201, 204})
	return err
}

func (s *Session) res2Err(res *http.Response, success []int) (err error) {
	for _, v := range success {
		if v == res.StatusCode {
			return nil
		}
	}
	return fmt.Errorf("%d %s", res.StatusCode, res.Status)
}

func (s *Session) doRequest(req *http.Request) (res *http.Response, err error) {
	res, err = s.client.Do(req)
	return
}

func (s *Session) abs(name string) (res string) {
	u, _ := url.Parse(s.base.String())
	u.Path = path.Join(u.Path, name)
	return u.String()
}
