package psui

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
	Auth    *BasicAuth `json:"auth"`
	PFTable string     `json:"table"`
	Devices []Device   `json:"devices"`
}

func (cfg *Config) Load(file string) error {
	contents, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("can't read config file: %w", err)
	}
	fmt.Printf("%s\n", contents)
	//x := Config{}
	return json.Unmarshal(contents, cfg)
	// fmt.Printf("%+v\n", x)
	// return err
}
