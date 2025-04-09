package config

type Config struct {
	config map[string]interface{}
}

// Get retrieves a value from the config map by key.
func (c *Config) Get(key string) interface{} {
	return c.config[key]
}

// Set allows setting a value in the config map.
func (c *Config) Set(key string, value interface{}) {
	c.config[key] = value
}

// Stage returns the current stage of the application.
//
// It is determined from APP_ENV variable
func (c *Config) Stage() string {
	return c.Get("APP_ENV").(string)
}

func NewConfig() *Config {
	return &Config{
		config: make(map[string]interface{}),
	}
}
