package crowdin

import (
	"io"
	"net/http"
	"time"
	"github.com/mreiferson/go-httpclient"
)

var (
	apiBaseURL    = "https://api.crowdin.com/api/project/"
)

type Crowdin struct {
	config struct {
			   apiBaseURL    string
			   token         string
			   client        *http.Client
		   }
	debug     bool
	logWriter io.Writer
}

func New(token string) *Crowdin {

	transport := &httpclient.Transport{
		ConnectTimeout:   5 * time.Second,
		ReadWriteTimeout: 40 * time.Second,
	}
	defer transport.Close()

	s := &Crowdin{}
	s.config.apiBaseURL = apiBaseURL
	s.config.token = token
	s.config.client = &http.Client{
		Transport: transport,
	}
	return s
}

func (crowdin *Crowdin) GetLanguageStatus(languageCode string) []File {
	// TODO
	return nil
}