package utils

import (
	"fmt"
	"testing"

	"github.com/wgliang/logcool/input/stdin"
)

func Test_RegistInputHandler(t *testing.T) {
	RegistInputHandler("stdin", stdininput.InitHandler)
}

func Test_RunInputs(t *testing.T) {
	config, err := LoadFromString(`
	{
		"input": [{
			"type": "file",
			"path": "./tmp/log/log.log",
			"sincedb_path": "",
			"start_position": "beginning"
		}]
	}
	`)
	RegistInputHandler("stdin", stdininput.InitHandler)
	if err != nil {
		fmt.Println(err)
	}
	config.RunInputs()
}
