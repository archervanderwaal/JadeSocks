package config

import (
	"errors"
	"github.com/BurntSushi/toml"
)

const (
	defaultListenAddr = ":8989"
)

type Config struct {
	ListenAddr string            `toml:"listen"`
	Users      map[string]string `toml:"users"`
}

func (conf *Config) LoadConfig(path string) error {
	md, err := toml.DecodeFile(path, conf)
	if err != nil {
		return err
	}
	if len(md.Undecoded()) > 0 {
		return errors.New("Unknown config keys in " + path)
	}
	if len(conf.ListenAddr) == 0 {
		conf.ListenAddr = defaultListenAddr
	}
	return nil
}