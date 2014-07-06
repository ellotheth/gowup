package gowup

import (
	"encoding/json"
	"net/url"
)

type JobSummary struct {
	Url Url    `json:"url"`
	Ip  string `json:"ip"`
}

type Url url.URL

func (u *Url) UnmarshalJSON(data []byte) error {
	var raw string

	// convert the json string to a real string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// parse the string into a Url
	parsed, err := url.Parse(raw)
	if err != nil {
		return err
	}
	*u = Url(*parsed)

	return nil
}
