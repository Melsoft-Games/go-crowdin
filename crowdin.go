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
	apiBaseURL        = "https://api.crowdin.com/api/project/"
	apiAccountBaseURL = "https://api.crowdin.com/api/account/"
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

// SetProject set project details
func (crowdin *Crowdin) SetProject(token, project string) *Crowdin {
	crowdin.config.token = token
	crowdin.config.project = project
	return crowdin
}

// SetDebug - traces errors if it's set to true.
func (crowdin *Crowdin) SetDebug(debug bool, logWriter io.Writer) {
	crowdin.debug = debug
	crowdin.logWriter = logWriter
}

// SetClient sets a custom http client. Can be useful in App Engine case.
func (crowdin *Crowdin) SetClient(client *http.Client) {
	crowdin.config.client = client
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

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiBaseURL+"%v/add-file?key=%v", crowdin.config.project, crowdin.config.token),
		params: params,
		files:  files,
	})

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

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiBaseURL+"%v/update-file?key=%v", crowdin.config.project, crowdin.config.token),
		params: params,
		files:  files,
	})

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

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiBaseURL+"%v/delete-file?key=%v", crowdin.config.project, crowdin.config.token),
		params: params,
	})

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

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiBaseURL+"%v/upload-translation?key=%v", crowdin.config.project, crowdin.config.token),
		params: params,
		files:  files,
	})

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

// GetTranslationsStatus - Track overall translation and proofreading progresses of each target language
func (crowdin *Crowdin) GetTranslationsStatus() ([]TranslationStatus, error) {

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiBaseURL+"%v/status?key=%v", crowdin.config.project, crowdin.config.token),
		params: map[string]string{
			"json": "",
		},
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	var responseAPI []TranslationStatus
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return responseAPI, nil
}

// GetLanguageStatus - Get the detailed translation progress for specified language.
// Language codes - https://crowdin.com/page/api/language-codes
func (crowdin *Crowdin) GetLanguageStatus(languageCode string) (*responseLanguageStatus, error) {

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiBaseURL+"%v/language-status?key=%v", crowdin.config.project, crowdin.config.token),
		params: map[string]string{
			"language": languageCode,
			"json":     "",
		},
	})

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

// GetProjectDetails - Get Crowdin Project details
func (crowdin *Crowdin) GetProjectDetails() (*ProjectInfo, error) {

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiBaseURL+"%v/info?key=%v", crowdin.config.project, crowdin.config.token),
		params: map[string]string{
			"json": "",
		},
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	var responseAPI ProjectInfo
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

//// DownloadTranslations - Download ZIP file with translations. You can choose the language of translation you need or download all of them at once.
//func (crowdin *Crowdin) DownloadTranslations() (error) {
//	// TODO
//}
//
//// ExportFile - This method exports single translated files from Crowdin. Additionally, it can be applied to export XLIFF files for offline localization.
//func (crowdin *Crowdin) ExportFile() (error) {
//	// TODO
//}
//
//// ExportTranslations - Build ZIP archive with the latest translations. Please note that this method can be invoked only once per 30 minutes (there is no such restriction for organization plans). Also API call will be ignored if there were no changes in the project since previous export. You can see whether ZIP archive with latest translations was actually build by status attribute ("built" or "skipped") returned in response.
//func (crowdin *Crowdin) ExportTranslations() (error) {
//	// TODO
//}

// AccountProjects - Get Crowdin Project details.
func (crowdin *Crowdin) GetAccountProjects(accountKey, loginUsername string) (*AccountDetails, error) {

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiAccountBaseURL+"get-projects?account-key=%v", accountKey),
		params: map[string]string{
			"login": loginUsername,
			"json":  "",
		},
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	var responseAPI AccountDetails
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// CreateProject - Create Crowdin project.
func (crowdin *Crowdin) CreateProject(accountKey, loginUsername string, options *CreateProjectOptions) (*responseManageProject, error) {

	params := make(map[string]string)
	params["json"] = ""
	params["login"] = loginUsername

	paramsArray := make(map[string][]string)

	if options != nil {

		if options.Name != "" {
			params["name"] = options.Name
		}

		if options.Identifier != "" {
			params["identifier"] = options.Identifier
		}

		if options.SourceLanguage != "" {
			params["source_language"] = options.SourceLanguage
		}

		if options.JoinPolicy != "" {
			params["join_policy"] = options.JoinPolicy
		}

		if options.Languages != nil {
			paramsArray["languages[]"] = options.Languages
		}
	}

	response, err := crowdin.post(&postOptions{
		urlStr:      fmt.Sprintf(apiAccountBaseURL+"create-project?account-key=%v", accountKey),
		params:      params,
		paramsArray: paramsArray,
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	var responseAPI responseManageProject
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// EditProject - Edit Crowdin project.
func (crowdin *Crowdin) EditProject(options *EditProjectOptions) (*responseManageProject, error) {

	params := make(map[string]string)
	params["json"] = ""

	paramsArray := make(map[string][]string)

	if options != nil {

		if options.Name != "" {
			params["name"] = options.Name
		}

		if options.JoinPolicy != "" {
			params["join_policy"] = options.JoinPolicy
		}

		if options.Languages != nil {
			paramsArray["languages[]"] = options.Languages
		}
	}

	response, err := crowdin.post(&postOptions{
		urlStr:      fmt.Sprintf(apiBaseURL+"%v/edit-project?key=%v", crowdin.config.project, crowdin.config.token),
		params:      params,
		paramsArray: paramsArray,
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	var responseAPI responseManageProject
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// DeleteProject - Delete Crowdin project with all translations.
func (crowdin *Crowdin) DeleteProject() (*responseDeleteProject, error) {

	params := make(map[string]string)
	params["json"] = ""

	response, err := crowdin.post(&postOptions{
		urlStr:      fmt.Sprintf(apiBaseURL+"%v/delete-project?key=%v", crowdin.config.project, crowdin.config.token),
		params:      params,
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	var responseAPI responseDeleteProject
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// AddDirectory - Add directory to Crowdin project.
// name - Directory name (with path if nested directory should be created).
func (crowdin *Crowdin) AddDirectory(directoryName string) (*responseGeneral, error) {

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiBaseURL+"%v/add-directory?key=%v", crowdin.config.project, crowdin.config.token),
		params: map[string]string{
			"name": directoryName,
			"json": "",
		},
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	var responseAPI responseGeneral
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// ChangeDirectory - Rename directory or modify its attributes. When renaming directory the path can not be changed (it means new_name parameter can not contain path, name only).
func (crowdin *Crowdin) ChangeDirectory(options *ChangeDirectoryOptions) (*responseGeneral, error) {

	params := make(map[string]string)
	params["json"] = ""

	if options != nil {

		if options.Name != "" {
			params["name"] = options.Name
		}

		if options.NewName != "" {
			params["new_name"] = options.NewName
		}

		if options.Title != "" {
			params["title"] = options.Title
		}
	}

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiBaseURL+"%v/change-directory?key=%v", crowdin.config.project, crowdin.config.token),
		params: params,
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	var responseAPI responseGeneral
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

// DeleteDirectory - Delete Crowdin project directory. All nested files and directories will be deleted too.
// name - Directory name (with path if nested directory should be created).
func (crowdin *Crowdin) DeleteDirectory(directoryName string) (*responseGeneral, error) {

	response, err := crowdin.post(&postOptions{
		urlStr: fmt.Sprintf(apiBaseURL+"%v/delete-directory?key=%v", crowdin.config.project, crowdin.config.token),
		params: map[string]string{
			"name": directoryName,
			"json": "",
		},
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	var responseAPI responseGeneral
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		log.Println(string(response))
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}

//// DownloadGlossary - Download Crowdin project glossaries as TBX file.
//func (crowdin *Crowdin) DownloadGlossary() (error) {
//	// TODO
//}
//
//// UploadGlossary - Upload your glossaries for Crowdin Project in TBX, CSV or XLS/XLSX file formats.
//func (crowdin *Crowdin) UploadGlossary() (error) {
//	// TODO
//}
//
//// SupportedLanguages - Get supported languages list with Crowdin codes mapped to locale name and standardized codes.
//func (crowdin *Crowdin) SupportedLanguages() (error) {
//	// TODO
//}
