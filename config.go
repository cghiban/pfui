package pfui

import (
	"encoding/json"
	"fmt"
	"os"
)

type Device struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
}

type BasicAuth struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Config struct {
	PFTable string     `json:"table"`
	Auth    *BasicAuth `json:"auth"`
	Devices []Device   `json:"devices"`
}

func (cfg *Config) Load(file string) error {
	contents, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("can't read config file: %w", err)
	}
	//fmt.Printf("%s\n", contents)
	return json.Unmarshal(contents, cfg)
}
