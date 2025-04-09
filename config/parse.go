package config

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// readYAMLFile reads and unmarshals a YAML file into the provided target interface.
func readYAMLFile(fileSys fs.FS, path string, target interface{}) error {
	yamlFile, err := fs.ReadFile(fileSys, path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(yamlFile, target); err != nil {
		return fmt.Errorf("error unmarshalling file '%s': %w", path, err)
	}
	return nil
}

// fileName extracts the file name without the extension.
func fileName(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// parseEnvVariables replaces environment variable placeholders in the config map.
func parseEnvVariables(config map[string]interface{}) map[string]interface{} {
	for key, value := range config {
		switch v := value.(type) {
		case map[string]interface{}:
			config[key] = parseEnvVariables(v)
		case string:
			if strings.HasPrefix(v, "$") {
				if val, ok := os.LookupEnv(strings.TrimLeft(v, "$")); ok {
					config[key] = val
				}
			}
		}
	}
	return config
}

// LoadConfigFile loads a YAML config file into a map[string]interface{}.
func LoadConfigFile(fileSys fs.FS, path string) (*Config, error) {
	var configs map[string]interface{}
	if err := readYAMLFile(fileSys, path, &configs); err != nil {
		return nil, err
	}
	config := NewConfig()
	config.config = parseEnvVariables(configs)
	return config, nil
}

// LoadDirectory loads a directory of YAML files into a map of map[string]interface{}.
func LoadDirectory(fileSys fs.FS, dirPath string) (map[string]*Config, error) {
	configs := make(map[string]*Config)
	files, err := fs.ReadDir(fileSys, dirPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".yaml" {
			continue
		}
		filePath := path.Join(dirPath, file.Name())
		data, err := LoadConfigFile(fileSys, filePath)
		if err != nil {
			log.Printf("Skipping file '%s': %v", file.Name(), err)
			continue
		}
		configs[fileName(file.Name())] = data
	}
	return configs, nil
}

// Merge merges multiple Config instances into one. It does not modify the original configs and returns the merged object
// with the last config's values taking precedence.
func Merge(configs ...*Config) *Config {
	merged := NewConfig()
	for _, config := range configs {
		for key, value := range config.config {
			merged.Set(key, value)
		}
	}
	return merged
}
