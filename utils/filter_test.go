package utils

import (
	"fmt"
	"testing"

	"github.com/wgliang/logcool/filter/zeus"
)

func Test_RegistFilterHandler(t *testing.T) {
	RegistFilterHandler("zeus", zeus.InitHandler)
}

func Test_RunFilters(t *testing.T) {
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
	RegistFilterHandler("zeus", zeus.InitHandler)
	if err != nil {
		fmt.Println(err)
	}
	config.RunFilters()
}
