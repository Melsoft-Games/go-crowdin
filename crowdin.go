package crowdin

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mreiferson/go-httpclient"
)

var (
	apiBaseURL = "https://api.crowdin.com/api/project/"
)

// Crowdin API wrapper
type Crowdin struct {
	config struct {
		apiBaseURL string
		token      string
		project    string
		client     *http.Client
	}
	debug     bool
	logWriter io.Writer
}

// New - create new instance of Crowdin API.
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

// SetDebug - traces errors if it's set to true.
func (crowdin *Crowdin) SetDebug(debug bool, logWriter io.Writer) {
	crowdin.debug = debug
	crowdin.logWriter = logWriter
}

// GetLanguageStatus - Get the detailed translation progress for specified language.
// Language codes - https://crowdin.com/page/api/language-codes
func (crowdin *Crowdin) GetLanguageStatus(languageCode string) (*responseLanguageStatus, error) {

	response, err := crowdin.post(fmt.Sprintf(apiBaseURL+"%v/language-status?key=%v", crowdin.config.project, crowdin.config.token),
		map[string]string{
			"language": languageCode,
			"json":     "",
		}, nil)

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	var responseAPI responseLanguageStatus
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// AddFile - Add new file to Crowdin project.
func (crowdin *Crowdin) AddFile(options *AddFileOptions) (*responseAddFile, error) {

	params := make(map[string]string)
	params["json"] = ""

	if options != nil {

		if options.Type != "" {
			params["type"] = options.Type
		}

		if options.Scheme != "" {
			params["scheme"] = options.Scheme
		}

		if options.FirstLineContainsHeader {
			params["first_line_contains_header"] = "true"
		} else {
			params["first_line_contains_header"] = "false"
		}

	}

	files := make(map[string]string)
	if options != nil && options.Files != nil {
		for k, path := range options.Files {
			files[fmt.Sprintf("files[%v]", k)] = path
		}
	}

	response, err := crowdin.post(fmt.Sprintf(apiBaseURL+"%v/add-file?key=%v", crowdin.config.project, crowdin.config.token),
		params,
		files)

	if err != nil {
		log.Println(string(response))
		return nil, err
	}

	var responseAPI responseAddFile
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil

}

// UpdateFile - Upload latest version of your localization file to Crowdin
func (crowdin *Crowdin) UpdateFile(options *UpdateFileOptions) (*responseGeneral, error) {

	params := make(map[string]string)
	params["json"] = ""

	if options != nil {

		if options.Scheme != "" {
			params["scheme"] = options.Scheme
		}

		if options.FirstLineContainsHeader {
			params["first_line_contains_header"] = "true"
		} else {
			params["first_line_contains_header"] = "false"
		}

	}

	files := make(map[string]string)
	if options != nil && options.Files != nil {
		for k, path := range options.Files {
			files[fmt.Sprintf("files[%v]", k)] = path
		}
	}

	response, err := crowdin.post(fmt.Sprintf(apiBaseURL+"%v/update-file?key=%v", crowdin.config.project, crowdin.config.token),
		params,
		files)

	if err != nil {
		log.Println(string(response))
		return nil, err
	}

	var responseAPI responseGeneral
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil

}

// DeleteFile - Delete file from Crowdin project. All the translations will be lost without ability to restore them
func (crowdin *Crowdin) DeleteFile(fileName string) (*responseGeneral, error) {

	params := make(map[string]string)
	params["json"] = ""
	params["file"] = fileName

	response, err := crowdin.post(fmt.Sprintf(apiBaseURL+"%v/delete-file?key=%v", crowdin.config.project, crowdin.config.token),
		params,
		nil)

	if err != nil {
		log.Println(string(response))
		return nil, err
	}

	var responseAPI responseGeneral
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil

}

// UpdateFile - Upload latest version of your localization file to Crowdin
func (crowdin *Crowdin) UploadTranslations(options *UploadTranslationsOptions) (*responseUploadTranslation, error) {

	params := make(map[string]string)
	params["json"] = ""

	if options != nil {

		if options.Language != "" {
			params["language"] = options.Language
		}

		params["import_duplicates"] = options.ImportDuplicates

	}

	files := make(map[string]string)
	if options != nil && options.Files != nil {
		for k, path := range options.Files {
			files[fmt.Sprintf("files[%v]", k)] = path
		}
	}

	response, err := crowdin.post(fmt.Sprintf(apiBaseURL+"%v/upload-translation?key=%v", crowdin.config.project, crowdin.config.token),
		params,
		files)

	if err != nil {
		log.Println(string(response))
		return nil, err
	}

	var responseAPI responseUploadTranslation
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil

}
