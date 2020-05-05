package config

import (
	"os"
	"testing"
)

func removeConfigFile() {
	_ = os.Remove(configFile(false))
}

func TestLoadConfig(t *testing.T) {
	removeConfigFile()
	config := LoadConfig(false)

	if config.ListenAddr != "" || config.RemoteAddr != "" || config.Password != "" {
		t.Error()
	}
}