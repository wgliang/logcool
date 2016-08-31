package utils

import (
	"fmt"
	"testing"
)

func Test_RegistFilterHandler(t *testing.T) {
	// RegistInputHandler("test")
}

func Test_RunFilters(t *testing.T) {
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
	config.RunFilters()
}
