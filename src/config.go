package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"go.yaml.in/yaml/v4"
)

var (
	configPath = expandTilde(path.Join("~", ".config", "poecampain", "config.yaml"))
)

func defaultClientPath() string {
	switch runtime.GOOS {
	case "windows":
		return `C:\Program Files (x86)\Steam\steamapps\common\Path of Exile\logs\Client.txt`
	case "linux":
		return expandTilde(`~/.steam/steam/steamapps/common/Path of Exile/logs/Client.txt`)
	case "darwin":
		return expandTilde(`~/Library/Application Support/Steam/steamapps/common/Path of Exile/logs/Client.txt`)
	default:
		return ""
	}
}

type Config struct {
	// PoE Client.txt
	Client string
}

func NewConfig() *Config {
	return &Config{
		Client: defaultClientPath(),
	}
}

func expandTilde(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return strings.Replace(path, "~", home, 1)
}

func readConfig() (*Config, error) {
	config := NewConfig()
	f, err := os.Open(configPath)
	if os.IsNotExist(err) {
		return config, nil
	}

	var c Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&c); err != nil {
		return nil, fmt.Errorf("failed to decode config: %v", err)
	}

	clientPath := expandTilde(c.Client)
	if _, err = os.Stat(clientPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("invalid client path in config: %v", err)
	}
	config.Client = clientPath

	return config, nil
}
