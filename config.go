package main
import (
	"encoding/json"
)
type Configuration struct {
	Listen string `json:"listen"`
	RemoteUrl string `json:"remote_url"`
	RemotePort int `json:"remote_port"`
}
func LoadConfig() (Configuration,error) {
	var conf Configuration
	json_string,err := fileGetContents("./env.json")
	if err != nil {
		return conf,err
	}
	err = json.Unmarshal(json_string, &conf)
	return conf,err
}
