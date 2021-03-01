package config

import (
	"encoding/json"
)

type Config struct {
	waa string `json: ""`
}

func Configure(data []byte) (Config, error) {
	c := Config{}

	err := json.Unmarshal([]byte(data), &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
