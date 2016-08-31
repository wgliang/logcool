package utils

import (
	"fmt"
	"testing"
)

func Test_RegistInputHandler(t *testing.T) {
	// RegistInputHandler("test")
}

func Test_RunInputs(t *testing.T) {
	config, err := LoadFromString(`
	{
		"input": [{
			"type": "file",
			"path": "/tmp/log/syslog",
			"sincedb_path": "",
			"start_position": "beginning"
		}]
	}
	`)
	if err != nil {
		fmt.Println(err)
	}
	config.RunInputs()
}
