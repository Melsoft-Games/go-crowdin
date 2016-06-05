package crowdin

import (
	"io"
	"net/http"
	"time"
	"github.com/mreiferson/go-httpclient"
	"fmt"
	"encoding/json"
	"log"
)

var (
	apiBaseURL = "https://api.crowdin.com/api/project/"
)

type Crowdin struct {
	config    struct {
				  apiBaseURL string
				  token      string
				  project    string
				  client     *http.Client
			  }
	debug     bool
	logWriter io.Writer
}

func New(token, project string) *Crowdin {

	transport := &httpclient.Transport{
		ConnectTimeout:   5 * time.Second,
		ReadWriteTimeout: 40 * time.Second,
	}
	defer transport.Close()

	s := &Crowdin{}
	s.config.apiBaseURL = apiBaseURL
	s.config.token = token
	s.config.project = project
	s.config.client = &http.Client{
		Transport: transport,
	}
	return s
}

// SetDebug traces errors if it's set to true.
func (crowdin *Crowdin) SetDebug(debug bool, logWriter io.Writer) {
	crowdin.debug = debug
	crowdin.logWriter = logWriter
}

func (crowdin *Crowdin) GetLanguageStatus(languageCode string) (*Files, error) {
	var files Files

	response, err := crowdin.post(fmt.Sprintf(apiBaseURL + "%v/language-status?key=%v", crowdin.config.project, crowdin.config.token),
		map[string]string{
			"language" : languageCode,
			"json" : "",
		}, nil)

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	err = json.Unmarshal(response, &files)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return &files, nil
}