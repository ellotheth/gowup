package gowup

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	apiEntryPoint string = "https://api.wheresitup.com/v4"
)

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

type Location struct {
	Name      string `json:"name"`
	Title     string `json:"title"`
	City      string `json:"location"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Lat       string `json:"latitude"`
	Lon       string `json:"longitutde"`
	Continent string `json:"continent_name"`
}

type WIU struct {
	Client string
	Token  string
}

func (api WIU) Locations() ([]Location, error) {
	response, err := api.get("sources")
	if err != nil {
		return nil, err
	}

	var body map[string][]Location
	if err := api.parse(response, &body); err != nil {
		return nil, err
	}

	sources, ok := body["sources"]
	if !ok {
		return nil, &Error{msg: "Locations missing from response"}
	}

	return sources, nil
}

func (api WIU) Jobs() (map[string]JobSummary, error) {
	response, err := api.get("jobs")
	if err != nil {
		return nil, err
	}

	var jobs map[string]JobSummary
	if err := api.parse(response, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (api WIU) setHeaders(req *http.Request, headers map[string]string) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth", "Bearer "+api.Client+" "+api.Token)

	for header, content := range headers {
		req.Header.Add(header, content)
	}
}

func (api WIU) parse(response *http.Response, body interface{}) error {
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err := json.Unmarshal(raw, body); err != nil {
		return err
	}

	return nil
}

func (api WIU) get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", apiEntryPoint+"/"+endpoint, nil)
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

	req, err := http.NewRequest("POST", apiEntryPoint+"/"+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	api.setHeaders(req, nil)

	return http.DefaultClient.Do(req)
}
