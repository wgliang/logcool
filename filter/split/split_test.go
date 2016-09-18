package split

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/wgliang/logcool/utils"
)

func init() {
	utils.RegistFilterHandler(ModuleName, InitHandler)
}

func Test_InitHandler(t *testing.T) {
	config := utils.ConfigRaw{}
	co, err := InitHandler(&config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(co)
}

func Test_Event(t *testing.T) {
	conf, err := utils.LoadFromString(`{
		"filter": [{
            "type": "split",
            "separator":";"
		}]
	}`)
	if err != nil {
		fmt.Println(err)
	}

	timestamp := time.Now()

	inchan := conf.Get(reflect.TypeOf(make(utils.InChan))).
		Interface().(utils.InChan)

	outchan := conf.Get(reflect.TypeOf(make(utils.OutChan))).
		Interface().(utils.OutChan)

	err = conf.RunFilters()

	inchan <- utils.LogEvent{
		Timestamp: timestamp,
		Message:   "filter test message",
	}

	event := <-outchan
	fmt.Println(event)
}
