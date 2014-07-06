package gowup

import (
	"encoding/json"
	"net/url"
	"time"
)

type JobSummary struct {
	Url       Url       `json:"url"`
	Ip        string    `json:"ip"`
	StartTime Time      `json:"start_time"`
}

type Url url.URL

type Time time.Time

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

// time.Time has an unmarshaler, but it assumes the JSON is coming in as a
// string. ours is coming in as a float64.
func (t *Time) UnmarshalJSON(data []byte) error {
	var raw float64

	// convert to a raw float
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// parse the string into a Url
	*t = Time(time.Unix(int64(raw), 0))

	return nil
}
