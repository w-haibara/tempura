package config

import (
	"encoding/json"

	"github.com/k0kubun/pp"
)

type Config struct {
	OutFile  string    `json: "outfile"`
	Greeting string    `json: "greeting"`
	Style    []string  `json: "style"`
	Commands []Command `json: "commands"`
}

type Command struct {
	Command string   `json: "command"`
	Message string   `json: "message"`
	Prompts []Prompt `json: "prompts"`
	Api     Api      `json: "api"`
}

type Prompt struct {
	Prompt string `json: "prompt"`
	Json   string `json: "json"`
}

type Api struct {
	Method string `json: "method"`
	Url    string `json: "url"`
}

const (
	GET    string = "GET"
	POST   string = "POST"
	PUT    string = "PUT"
	DELETE string = "DELETE"
)

func Configure(data []byte) (Config, error) {
	c := Config{}

	err := json.Unmarshal([]byte(data), &c)
	if err != nil {
		return Config{}, err
	}

	pp.Println(c)

	return c, nil
}
