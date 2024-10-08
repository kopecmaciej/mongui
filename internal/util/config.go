package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

const (
	ConfigDir = "vi-mongo"
)

// MergeConfigs merges the loaded config with the default config
func MergeConfigs(loaded, defaultConfig interface{}) {
	mergeConfigsRecursive(reflect.ValueOf(loaded).Elem(), reflect.ValueOf(defaultConfig).Elem())
}

// mergeConfigsRecursive recursively merges nested structs
func mergeConfigsRecursive(loaded, defaultValue reflect.Value) {
	for i := 0; i < loaded.NumField(); i++ {
		field := loaded.Field(i)
		defaultField := defaultValue.Field(i)

		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				field.Set(defaultField)
			}
		case reflect.Slice:
			if field.Len() == 0 {
				field.Set(defaultField)
			}
		case reflect.Struct:
			mergeConfigsRecursive(field, defaultField)
		}
	}
}

// LoadConfigFile loads a configuration file, merges it with defaults, and returns the result
func LoadConfigFile[T any](defaultConfig *T, configPath string) (*T, error) {
	// Ensure the config directory exists
	err := ensureConfigDirExist()
	if err != nil {
		return nil, err
	}

	// Read the config file
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, create it with default settings
			bytes, err = marshalConfig(defaultConfig, configPath)
			if err != nil {
				return nil, err
			}
			err = os.WriteFile(configPath, bytes, 0644)
			if err != nil {
				return nil, err
			}
			return defaultConfig, nil
		}
		return nil, err
	}

	// Unmarshal the config file
	config := new(T)
	err = unmarshalConfig(bytes, configPath, config)
	if err != nil {
		return nil, err
	}

	// Merge loaded config with default config
	MergeConfigs(config, defaultConfig)

	return config, nil
}

// marshalConfig marshals the config based on the file extension
func marshalConfig[T any](config *T, configPath string) ([]byte, error) {
	switch filepath.Ext(configPath) {
	case ".json":
		return json.Marshal(config)
	case ".yaml", ".yml":
		return yaml.Marshal(config)
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", configPath)
	}
}

// unmarshalConfig unmarshals the config based on the file extension
func unmarshalConfig[T any](data []byte, configPath string, config *T) error {
	switch filepath.Ext(configPath) {
	case ".json":
		return json.Unmarshal(data, config)
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, config)
	default:
		return fmt.Errorf("unsupported file extension: %s", configPath)
	}
}

// ensureConfigDirExist ensures the config directory exists
// If it does not exist, it will be created
func ensureConfigDirExist() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return os.MkdirAll(configDir, 0755)
	}
	return nil
}

// GetConfigDir returns the path to the config directory
func GetConfigDir() (string, error) {
	configPath, err := xdg.ConfigFile(ConfigDir)
	if err != nil {
		return "", err
	}
	return configPath, nil
}
