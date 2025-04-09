package config

import (
	"fmt"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestReadYAMLFile(t *testing.T) {
	fs := fstest.MapFS{
		"test.yaml": {Data: []byte("key: value")},
	}

	var result map[string]interface{}
	err := readYAMLFile(fs, "test.yaml", &result)
	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])
}

func TestFileName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"config.yaml", "config"},
		{"config.json", "config"},
		{"config", "config"},
	}

	for _, test := range tests {
		result := fileName(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestParseEnvVariables(t *testing.T) {
	os.Setenv("TEST_ENV", "env_value")
	defer os.Unsetenv("TEST_ENV")

	input := map[string]interface{}{
		"key1": "$TEST_ENV",
		"key2": "static_value",
		"nested": map[string]interface{}{
			"key3": "$TEST_ENV",
		},
	}

	expected := map[string]interface{}{
		"key1": "env_value",
		"key2": "static_value",
		"nested": map[string]interface{}{
			"key3": "env_value",
		},
	}

	result := parseEnvVariables(input)
	assert.Equal(t, expected, result)
}

func TestLoadConfigFile(t *testing.T) {
	fs := fstest.MapFS{
		"test.yaml": {Data: []byte("key: value")},
	}

	config, err := LoadConfigFile(fs, "test.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "value", config.config["key"])
}

func TestLoadDirectory(t *testing.T) {
	fs := fstest.MapFS{
		"config1.yaml": {Data: []byte("key1: value1")},
		"config2.yaml": {Data: []byte("key2: value2")},
	}

	configs, err := LoadDirectory(fs, ".")
	fmt.Printf("confnigd: %#v", configs)
	assert.NoError(t, err)
	assert.Equal(t, "value1", configs["config1"].config["key1"])
	assert.Equal(t, "value2", configs["config2"].config["key2"])
}
