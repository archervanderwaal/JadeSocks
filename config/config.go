package config

import (
	"fmt"
	"github.com/archervanderwaal/JadeSocks/utils"
	"github.com/aybabtme/rgbterm"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	Listen                         = "listen"
	Remote                         = "remote"
	Password                       = "password"
	configFileName                 = "JadeSocks.yaml"
	serverModeDefaultConfigContent = `
listen: :5678
password: admin
`
	clientModeDefaultConfigContent = `
listen: :5678
remote: x.x.x.x:5678
password: admin
`
)

type Config struct {
	ListenAddr string `yaml:"listen"`
	RemoteAddr string `yaml:"remote"`
	Password   string `yaml:"password"`
}

func LoadConfig(serverMode bool) *Config {
	var settings Config
	configFile, err := ioutil.ReadFile(configFile(serverMode))
	if err != nil {
		log.Println(rgbterm.FgString("Internal error "+err.Error(), 255, 0, 0))
		os.Exit(1)
	}
	_ = yaml.Unmarshal(configFile, &settings)
	return &settings
}

func configFile(serverMode bool) string {
	configFilePath := filepath.Join(utils.Home(), configFileName)
	writeDefaultConfigContent(serverMode)
	return configFilePath
}

func writeDefaultConfigContent(serverMode bool) {
	if !utils.Exists(utils.Home()) {
		_ = os.Mkdir(utils.Home(), 0755)
	}
	if utils.Exists(filepath.Join(utils.Home(), configFileName)) {
		return
	}
	file, err := os.Create(filepath.Join(utils.Home(), configFileName))
	if err != nil {
		log.Println(rgbterm.FgString("Internal error "+err.Error(), 255, 0, 0))
		os.Exit(1)
	}
	defer file.Close()
	if serverMode {
		_, _ = fmt.Fprint(file, serverModeDefaultConfigContent)
	} else {
		_, _ = fmt.Fprint(file, clientModeDefaultConfigContent)
	}
}
