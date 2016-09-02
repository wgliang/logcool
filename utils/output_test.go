package utils

import (
	"fmt"
	"testing"

	"github.com/wgliang/logcool/output/stdout"
)

func Test_RegistOutputHandler(t *testing.T) {
	RegistOutputHandler("stdout", outputstdout.InitHandler)
}

func Test_RunOutputs(t *testing.T) {
	config, err := LoadFromString(`
	{
		"input": [{
			"type": "file",
			"path": "/tmp/log/log.log",
			"sincedb_path": "",
			"start_position": "beginning"
		}]
	}
	`)
	RegistOutputHandler("stdout", outputstdout.InitHandler)
	if err != nil {
		fmt.Println(err)
	}
	config.RunOutputs()
}
