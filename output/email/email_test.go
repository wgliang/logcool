package outputemail

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/wgliang/logcool/utils"
)

func init() {
	utils.RegistOutputHandler(ModuleName, InitHandler)
}

func Test_InitHandler(t *testing.T) {
	conf, err := utils.LoadFromString(`{
		"output": [{
	           "type": "stdout"
	       }]
	}`)
	var confraw *utils.ConfigRaw
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}
	InitHandler(confraw)
}

func Test_Event(t *testing.T) {
	conf, err := utils.LoadFromString(`{
		"output": [{
			"type": "stdout"
		}]
	}`)

	err = conf.RunOutputs()
	if err != nil {
		fmt.Println(err)
	}
	outchan := conf.Get(reflect.TypeOf(make(utils.OutChan))).
		Interface().(utils.OutChan)
	outchan <- utils.LogEvent{
		Timestamp: time.Now(),
		Message:   "outputstdout test message",
	}

	time.Sleep(time.Duration(5) * time.Second)
}
