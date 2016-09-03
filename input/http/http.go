package httpinput

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/wgliang/logcool/utils"
	"github.com/wgliang/logcool/utils/logevent"
)

const (
	ModuleName = "http"
)

type InputConfig struct {
	config.InputConfig
	Addr      string   `json:"addr"`
	Method    []string `json:"method"`
	Urls      string   `json:"urls"`
	Intervals int      `json:"intervals"`

	hostname string
	httpChan chan logevent.LogEvent
}

func InitHandler(confraw *config.ConfigRaw) (retconf config.TypeInputConfig, err error) {
	conf := InputConfig{
		InputConfig: config.InputConfig{
			CommonConfig: config.CommonConfig{
				Type: ModuleName,
			},
		},
		Method:   "GET",
		Interval: 60,
	}
	if err = config.ReflectConfig(confraw, &conf); err != nil {
		return
	}

	if conf.hostname, err = os.Hostname(); err != nil {
		return
	}
	fmt.Println("=============")
	retconf = &conf
	return
}

func (t *InputConfig) Start() {
	fmt.Println("start http....")
	t.Invoke(t.listen)
}

func (ic *InputConfig) listen(logger *logrus.Logger, inchan config.InChan) {
	http.HandleFunc("/logcool", ic.Handler)
	//http server.
	http.ListenAndServe(ic.Addr, nil)
	fmt.Println(`Now Serving...`)

	for {
		select {
		case event := <-ic.httpChan:
			inchan <- event
		}
	}
}

// Handler 处理请求
func (ic *InputConfig) Handler(w http.ResponseWriter, r *http.Request) {
	var (
		message string
		ok      bool
	)

	if r.Method == "GET" {
		message, ok := r.Form["data"]
		if ok == false {
			message = "error"
		}

	} else if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		message = result
	}
	event := logevent.LogEvent{
		Timestamp: time.Now(),
		Message:   message,
		Extra: map[string]interface{}{
			"host": t.hostname,
		},
	}
	return
}
