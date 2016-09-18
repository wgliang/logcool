package outputexec

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
	           "type": "lexec"
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
			"type": "lexec"
		}]
	}`)

	err = conf.RunOutputs()
	if err != nil {
		fmt.Println(err)
	}
	outchan := conf.Get(reflect.TypeOf(make(utils.OutChan))).
		Interface().(utils.OutChan)
	args := make(map[string]interface{})
	args["args"] = []string{"-e"}
	outchan <- utils.LogEvent{
		Timestamp: time.Now(),
		Message:   "ps",
		Extra:    args,
	}

	time.Sleep(time.Duration(5) * time.Second)
}
