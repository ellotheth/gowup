package gowup

import (
	"encoding/json"
	"net/url"
	"time"
)

type JobSummary struct {
	Url        Url       `json:"url"`
	Ip         string    `json:"ip"`
	StartTime  Time      `json:"start_time"`
	ExpireTime Time      `json:"expiry"`
	Services   []Service `json:"services"`
}

type Service struct {
	Server string   `json:"server"`
	Tests  []string `json:"checks"`
}

type Url struct {
	*url.URL
}

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
	u.URL = parsed

	return nil
}

type Time struct {
	time.Time
}

// time.Time has an unmarshaler, but it assumes the JSON is coming in as a
// string. ours is coming in as a float64, and sometimes it comes in a hashmap.
func (t *Time) UnmarshalJSON(data []byte) error {
	var raw interface{}

	// convert to a raw float
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// this is just freaking black magic right here, i don't even know
	switch v := raw.(type) {
	case float64:
		t.Time = time.Unix(int64(v), 0)
	case map[string]interface{}:
		sec, ok := v["sec"].(float64)
		if !ok {
			return &Error{msg: "No seconds found in expiry time type"}
		}
		t.Time = time.Unix(int64(sec), 0)
	default:
		return &Error{msg: "I have no idea what to do with this time type you gave me"}
	}

	return nil
}
