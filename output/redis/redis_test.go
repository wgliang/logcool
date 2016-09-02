package outputredis

import (
	"fmt"
	// "reflect"
	"testing"
	// "time"

	"github.com/wgliang/logcool/utils"
	// "github.com/wgliang/logcool/utils/logevent"
)

func init() {
	utils.RegistOutputHandler(ModuleName, InitHandler)
}

func Test_InitHandler(t *testing.T) {
	conf, err := utils.LoadFromString(`{
		"output": [{
	           "type": "redis",
	           "key": "logcool",
	           "host": "127.0.0.1:6379",
	           "password":"",
	           "data_type": "list",
	           "timeout": 5,
	           "reconnect_interval": 1
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
			"type": "redis",
	           "key": "logcool",
	           "host": "127.0.0.1:6379",
	           "password":"",
	           "data_type": "list",
	           "timeout": 5,
	           "reconnect_interval": 1
		}]
	}`)

	err = conf.RunOutputs()
	if err != nil {
		fmt.Println(err)
	}

	// evchan := conf.Get(reflect.TypeOf(make(chan logevent.LogEvent))).
	// 	Interface().(chan logevent.LogEvent)

	// evchan <- logevent.LogEvent{
	// 	Timestamp: time.Now(),
	// 	Message:   "outputstdout test message",
	// }
}
