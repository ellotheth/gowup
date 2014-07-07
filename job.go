package gowup

import (
	"encoding/json"
	"net/url"
	"time"
)

type Job struct {
	Summary JobSummary `json:"request"`
	Details JobDetails `json:"response"`
}

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

type JobRequest struct {
	Url       string   `json:"uri"`
	Tests     []string `json:"tests"`
	Locations []string `json:"sources"`
}

type JobDetails struct {
	Done    JobDetail `json:"complete"`
	NotDone JobDetail `json:"in_progress"`
	Error   JobDetail `json:"error"`
}

type JobDetail map[string]map[string]interface{}

func (j *JobDetail) UnmarshalJSON(data []byte) error {
	var raw interface{}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case []interface{}:
		*j = JobDetail{}
	case map[string]interface{}:
		job := JobDetail{}

		// there must be a better way.
		for city, tests := range v {
			job[city] = map[string]interface{}{}
			for test, details := range tests.(map[string]interface{}) {
				content, ok := details.(map[string]interface{})["summary"]
				if ok {
					job[city][test] = content
				}
			}
		}
		*j = job

	}

	return nil
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
		// there's an easier way to do this with json flags. oops
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
