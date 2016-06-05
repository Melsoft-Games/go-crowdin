package crowdin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"log"
	"net/url"
)

func (crowdin *Crowdin) post(urlStr string, params map[string]string) ([]byte, error) {

	form := url.Values{}
	if params != nil {
		for k, v := range params {
			form.Add(k, v)
		}
	}

	response, err := crowdin.config.client.PostForm(urlStr, form)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = APIError{What: fmt.Sprintf("Status code: %v", response.StatusCode)}
		crowdin.log(err)
		return nil, err
	}

	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return bodyResponse, nil
}

func (crowdin *Crowdin) log(a interface{}) {
	if crowdin.debug {
		log.Println(a)
		if crowdin.logWriter != nil {
			timestamp := time.Now().Format(time.RFC3339)
			msg := fmt.Sprintf("%v: %v", timestamp, a)
			fmt.Fprintln(crowdin.logWriter, msg)
		}
	}
}

// APIError holds data of errors returned from the API.
type APIError struct {
	What string
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v", e.What)
}
