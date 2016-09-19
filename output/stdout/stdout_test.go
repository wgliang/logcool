package outputstdout

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
	        "type": "email",
            "server":"smtp.163.com:25",
            "from":"wgl672@163.com",
            "password":"wgl5124991",
            "to":["liangcszzu@163.com"],
            "cc":"liangcszzu@163.com"
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
			"type": "email",
            "server":"smtp.163.com:25",
            "from":"wgl672@163.com",
            "password":"wgl5124991",
            "to":["liangcszzu@163.com"],
            "cc":"liangcszzu@163.com"
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
