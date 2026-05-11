package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configDir = ".simplemdm-cli"
const configFile = "config.json"

type Context struct {
	Name string `json:"name,omitempty"`
}

type Config struct {
	DefaultContext string             `json:"default_context"`
	Contexts       map[string]Context `json:"contexts"`
}

func NewConfig() *Config {
	return &Config{
		Contexts: make(map[string]Context),
	}
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, configDir)
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), configFile)
}

func Load() (*Config, error) {
	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := NewConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}

func (c *Config) Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(ConfigPath(), data, 0600)
}

func (c *Config) SetContext(name string, ctx Context) {
	c.Contexts[name] = ctx
}

func (c *Config) DeleteContext(name string) {
	delete(c.Contexts, name)
}

func (c *Config) GetContext(name string) (Context, bool) {
	ctx, ok := c.Contexts[name]
	return ctx, ok
}

func (c *Config) GetDefaultContext() (string, Context, bool) {
	if c.DefaultContext == "" {
		return "", Context{}, false
	}
	ctx, ok := c.Contexts[c.DefaultContext]
	return c.DefaultContext, ctx, ok
}
