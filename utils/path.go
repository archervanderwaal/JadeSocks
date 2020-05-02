package utils

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

var homeDir = ""

const folder = ".JadeSocks"

func Home() string {
	if homeDir != "" {
		return homeDir
	}
	if h, err := homedir.Dir(); err == nil {
		homeDir = filepath.Join(h, folder)
	} else {
		cwd, err := os.Getwd()
		if err == nil {
			homeDir = filepath.Join(cwd, folder)
		} else {
			homeDir = folder
		}
	}
	return homeDir
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}