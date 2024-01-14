package config_helper

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config interface {
	Validate() error
}

var defaultConfigFileName = "config.yaml"

func RetrievePropertiesFromYaml(filePath string, cfg Config) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		return err
	}
	return nil
}

func LoadDefaultConfig(cfg Config) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, defaultConfigFileName)

	return RetrievePropertiesFromYaml(configPath, cfg)
}
