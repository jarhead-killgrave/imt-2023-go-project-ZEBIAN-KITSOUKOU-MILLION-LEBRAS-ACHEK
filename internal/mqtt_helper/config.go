package mqtt_helper

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
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
		logutil.Error("Failed to open file:\n\t << %v >>", err)
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var cfg MQTTConfig
	decoder := yaml.NewDecoder(f)
	if decoder.Decode(&cfg) != nil {
		logutil.Error("Failed to decode file: << %v >>", err)
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		logutil.Error("Failed to validate config: << %v >>", err)
		return nil, err
	}

	return &cfg, err
}

func (c *MQTTConfig) GetServerAddress() string {
	return fmt.Sprintf("%s://%s:%d", c.Server.Protocol, c.Server.Host, c.Server.Port)
}

func (c *MQTTConfig) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("host is empty")
	}

	if c.Server.Port == 0 {
		return fmt.Errorf("port is empty")
	}

	if c.Server.Protocol == "" {
		return fmt.Errorf("protocol is empty")
	}

	return nil
}