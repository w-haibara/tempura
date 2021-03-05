package config

import (
	"encoding/json"
	"reflect"

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
	Print   Print    `json: "print"`
}

type Prompt struct {
	Prompt string `json: "prompt"`
	Mask   bool   `json: mask`
	Json   string `json: "json"`
	Header string `json: "header"`
	Query  string `json: "query"`
}

type Api struct {
	Method string `json: "method"`
	Url    string `json: "url"`
}

type Print struct {
	Json   bool `json: "json"`
	Header bool `json: "header"`
	Query  bool `json: "query"`
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

func IsEmpty(p []Prompt, fname string) bool {
	isEmpty := true
	for i, _ := range p {
		isEmpty = isEmpty && (reflect.ValueOf(p[i]).FieldByName(fname).String() == "")
	}
	return isEmpty
}
