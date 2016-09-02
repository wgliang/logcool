package utils

import (
	"fmt"
	"testing"

	"github.com/codegangsta/inject"
)

func Test_SetInjector(t *testing.T) {
	comc := CommonConfig{
		Type: "test",
	}
	inj := inject.New()
	comc.SetInjector(inj)
	fmt.Println(comc)
}

func Test_GetType(t *testing.T) {
	comc := CommonConfig{
		Type: "test",
	}
	fmt.Println(comc.Type)
}

func Test_Invoke(t *testing.T) {
	conf, _ := LoadFromString(Defaultconfig)
	fmt.Println(conf)
	// inj := inject.New()
	// conf.Invoke(inj)
}

func Test_LoadFromFile(t *testing.T) {
	LoadFromFile("../templates/stdin2stdout.json")
}

func Test_LoadFromString(t *testing.T) {
	LoadFromString(`
	{
		"input": [{
			"type": "file",
			"path": "./tmp/log/log.log",
			"sincedb_path": "",
			"start_position": "beginning"
		}]
	}
	`)
}

func Test_LoadFromData(t *testing.T) {
	LoadFromData([]byte(`
	{
		"input": [{
			"type": "file",
			"path": "./tmp/log/log.log",
			"sincedb_path": "",
			"start_position": "beginning"
		}]
	}
	`))
}

func Test_ReflectConfig(t *testing.T) {
	// todo
}

func Test_CleanComments(t *testing.T) {
	if data, err := CleanComments([]byte(`
	{
		"input": [{
			"type": "file",
			"path": "/tmp/log/syslog",
			"sincedb_path": "",
			"start_position": "beginning"
		}]
	}
	`)); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(data)
	}
}

func Test_InvokeSimple(t *testing.T) {
	// todo
}
