package utils

import (
	"os"
	"testing"
)

func TestParseArgs(t *testing.T) {
	os.Args = []string{
		"JadeSocks",
		"-s",
		"~/.JadeSocks.json",
	}
	content, args := ParseArgs(os.Args)
	if len(content) == 1 && len(args) == 1 {
		return
	}
	if content[0] != "~/.JadeSocks.json" || args[0] != "-s" {
		t.Fail()
	}
}
