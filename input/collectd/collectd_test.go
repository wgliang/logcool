package collectdinput

import (
	"fmt"
	"testing"
	"time"

	"github.com/wgliang/logcool/utils"
)

func init() {
	utils.RegistInputHandler(ModuleName, InitHandler)
}

func Test_InitHandler(t *testing.T) {
	config := utils.ConfigRaw{}
	co, err := InitHandler(&config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(co)
}

func Test_Start(t *testing.T) {
	conf, err := utils.LoadFromString(`{
		"input": [{
			"type": "collectd"
		}]
	}`)
	err = conf.RunInputs()
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Duration(5) * time.Second)
}
