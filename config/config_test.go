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

	if config.ListenAddr != "localhost:1087" || config.RemoteAddr != "xxx.xxx.xxx.xxx:8989" {
		t.Error()
	}
}