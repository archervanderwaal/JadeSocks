package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf := &Config{}
	err := conf.LoadConfig("JadeSocks.toml")
	if err != nil {
		t.Fail()
	}
}