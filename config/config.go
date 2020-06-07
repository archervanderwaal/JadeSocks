package config

import (
	"fmt"
	"github.com/archervanderwaal/JadeSocks/logger"
	"github.com/archervanderwaal/JadeSocks/utils"
	"github.com/aybabtme/rgbterm"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	clientConfigFileName           = "JadeSocks-client.yaml"
	serverConfigFileName           = "JadeSocks-server.yaml"
	defaultClientConfigFileContent = `
listen: localhost:1087
remote: xxx.xxx.xxx.xxx:8989
users:
  user1: passwd2
`
	defaultServerConfigFileContent = `
listen: localhost:8989
users:
  user1: passwd1
  user2: passwd2
  user3: passwd3
`
)

type Config struct {
	ListenAddr string            `yaml:"listen"`
	RemoteAddr string            `yaml:"remote"`
	Users      map[string]string `yaml:"users"`
}

func LoadConfig(serverMode bool) *Config {
	var settings Config
	configFile, err := ioutil.ReadFile(configFile(serverMode))
	if err != nil {
		logger.Logger.Error(rgbterm.FgString("Internal error "+err.Error(), 255, 0, 0))
		os.Exit(1)
	}
	_ = yaml.Unmarshal(configFile, &settings)
	return &settings
}

func configFile(serverMode bool) string {
	var configFilePath string
	if serverMode {
		configFilePath = filepath.Join(utils.Home(), serverConfigFileName)
	} else {
		configFilePath = filepath.Join(utils.Home(), clientConfigFileName)
	}
	writeDefaultConfigContent(serverMode, configFilePath)
	return configFilePath
}

func writeDefaultConfigContent(serverMode bool, configFilePath string) {
	if !utils.Exists(utils.Home()) {
		_ = os.Mkdir(utils.Home(), 0755)
	}
	if utils.Exists(configFilePath) {
		return
	}
	file, err := os.Create(configFilePath)
	if err != nil {
		log.Println(rgbterm.FgString("Internal error "+err.Error(), 255, 0, 0))
		os.Exit(1)
	}
	defer file.Close()
	if serverMode {
		_, _ = fmt.Fprint(file, defaultServerConfigFileContent)
	} else {
		_, _ = fmt.Fprint(file, defaultClientConfigFileContent)
	}
}
