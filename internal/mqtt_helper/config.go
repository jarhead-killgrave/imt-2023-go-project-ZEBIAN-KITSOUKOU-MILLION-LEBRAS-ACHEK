package mqtt_helper

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type MQTTConfig struct {
	Server struct {
		Protocol string `yaml:"protocol"`
		Port     int    `yaml:"port"`
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"server"`
}

func RetrievePropertiesFromConfig(filePath string) (*MQTTConfig, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("\033[31mfailed to open file:\n\t<<%w>>\033[0m", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var cfg MQTTConfig
	decoder := yaml.NewDecoder(f)
	if decoder.Decode(&cfg) != nil {
		return nil, fmt.Errorf("\033[31mfailed to decode file:\n\t<<%w>>\033[0m", err)
	}
	return &cfg, err
}
