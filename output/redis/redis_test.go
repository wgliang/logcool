package outputredis

import (
	"fmt"
	"testing"
	"time"

	"github.com/wgliang/logcool/utils"
	"github.com/wgliang/logcool/utils/logevent"
)

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
	conf := OutputConfig{
		OutputConfig: utils.OutputConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
		Key:               "logcool",
		DataType:          "list",
		Timeout:           5,
		ReconnectInterval: 1,

		evchan: make(chan logevent.LogEvent),
	}
	var confraw *utils.ConfigRaw
	if err := utils.ReflectConfig(confraw, &conf); err != nil {
		fmt.Println(err)
		return
	}
	InitHandler(confraw)
	ev := logevent.LogEvent{
		Timestamp: time.Now(),
		Message:   "outputredis test message",
	}
	conf.Event(ev)
}
