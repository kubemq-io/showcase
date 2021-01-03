package config

import (
	"encoding/json"
	"io/ioutil"
)

type WebConfig struct {
	ApiServerURL string `json:"API_SERVER_URL"`
	PollInterval int    `json:"POLL_INTERVAL"`
}

func (w *WebConfig) Write(filename string) error {
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
