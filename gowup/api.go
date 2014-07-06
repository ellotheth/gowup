package gowup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	apiEntryPoint string = "https://api.wheresitup.com/v4"
)

type WIUError struct {
	msg string
}

func (e *WIUError) Error() string {
	return e.msg
}

type WIU struct {
	Client string
	Token  string
}

func (api WIU) Locations() ([]map[string]string, error) {
	response, err := api.get("sources")
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var body map[string][]map[string]string
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}

	sources, ok := body["sources"]
	if !ok {
		return nil, &WIUError{msg: "Locations missing from response"}
	}

	return sources, nil
}

func (api WIU) setHeaders(req *http.Request, headers map[string]string) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth", fmt.Sprintf("Bearer %s %s", api.Client, api.Token))

	for header, content := range headers {
		req.Header.Add(header, content)
	}
}

func (api WIU) get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", apiEntryPoint, endpoint), nil)
	if err != nil {
		return nil, err
	}

	api.setHeaders(req, nil)

	return http.DefaultClient.Do(req)
}

func (api WIU) post(endpoint string, data interface{}) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", apiEntryPoint, endpoint), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	api.setHeaders(req, nil)

	return http.DefaultClient.Do(req)
}
