package guppy

import (
	"fmt"
	"net/http"
)

var (
	apiEntryPoint string = "https://api.wheresitup.com/v4"
)

type WIU struct {
	Client string
	Token  string
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
