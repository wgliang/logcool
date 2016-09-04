package httpinput

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/wgliang/logcool/utils"
)

func httpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
}

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
			"type": "http",
            "addr": "127.0.0.1:6789",
            "urls": "/logcool",
            "method": ["HEAD","GET","POST"],
            "intervals":5
		}]
	}`)

	err = conf.RunInputs()
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Duration(1) * time.Second)
	httpGet("http://127.0.0.1:6789/logcool?data=logcool")
}
