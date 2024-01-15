package config

import "fmt"

type API struct {
	ServerEndpoint struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"api"`

	Database InfluxDBSettings `yaml:"database"`
}

func (c *API) Validate() error {
	if c.ServerEndpoint.Host == "" {
		return fmt.Errorf("host is empty")
	}

	if c.ServerEndpoint.Port == "" {
		return fmt.Errorf("port is empty")
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	return nil
}

func LoadDefaultAPIConfig() (*API, error) {
	var cfg API
	err := LoadDefaultConfig(&cfg)

	return &cfg, err
}
