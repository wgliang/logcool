package httpinput

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/wgliang/logcool/utils"
	"github.com/wgliang/logcool/utils/logevent"
)

const (
	ModuleName = "http"
)

type InputConfig struct {
	utils.InputConfig
	Addr      string   `json:"addr"`
	Method    []string `json:"method"`
	Urls      string   `json:"urls"`
	Intervals int      `json:"intervals"`

	hostname string
	httpChan chan logevent.LogEvent
}

func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeInputConfig, err error) {
	conf := InputConfig{
		InputConfig: utils.InputConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
		Method:    []string{"GET"},
		Intervals: 10,
	}
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}

	if conf.hostname, err = os.Hostname(); err != nil {
		return
	}
	fmt.Println("=============")
	retconf = &conf
	return
}

func (ic *InputConfig) Start() {
	fmt.Println("start http....")
	ic.Invoke(ic.listen)
}

func (ic *InputConfig) listen(logger *logrus.Logger, inchan utils.InChan) {
	var mux = http.NewServeMux()
	mux.HandleFunc(ic.Urls, ic.Handler)
	fmt.Println(ic.Addr)
	//http server.
	go func(serverAddr string, m *http.ServeMux) {
		if err := http.ListenAndServe(serverAddr, m); err != nil {
			fmt.Println(err)
		}
	}(ic.Addr, mux)

	fmt.Println(`Now Serving...`)

	for {
		select {
		case event := <-ic.httpChan:
			inchan <- event
		}
		time.Sleep(time.Second)
	}
}

// Handler 处理请求
func (ic *InputConfig) Handler(w http.ResponseWriter, r *http.Request) {
	var message string

	if r.Method == "GET" {
		// if _, ok := r.Form["data"]; ok {
		// 	if len(r.Form["data"]) > 0 {
		// 		message = r.Form["data"][0]
		// 	}
		// }
		message = "logcool"
	} else if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		message = string(result)
	}
	event := logevent.LogEvent{
		Timestamp: time.Now(),
		Message:   message,
		Extra: map[string]interface{}{
			"host": ic.hostname,
		},
	}
	ic.httpChan <- event
	w.Write([]byte(message))
	return
}
