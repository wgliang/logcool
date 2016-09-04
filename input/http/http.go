package httpinput

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/glog"

	"github.com/wgliang/logcool/utils"
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
	httpChan chan utils.LogEvent
}

func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeInputConfig, err error) {
	glog.Infoln(ModuleName + " input-plug Init...")
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
		glog.Errorln(err)
		return
	}

	if conf.hostname, err = os.Hostname(); err != nil {
		glog.Errorln(err)
		return
	}
	conf.httpChan = make(chan utils.LogEvent, 10)
	retconf = &conf
	return
}

func (ic *InputConfig) Start() {
	glog.Infoln(ModuleName + " input-plug Starting...")
	ic.Invoke(ic.listen)
}

func (ic *InputConfig) listen(logger *logrus.Logger, inchan utils.InChan) {
	var mux = http.NewServeMux()
	mux.HandleFunc(ic.Urls, ic.Handler)
	//http server.
	go func(serverAddr string, m *http.ServeMux) {
		if err := http.ListenAndServe(serverAddr, m); err != nil {
			glog.Errorln(err)
		}
	}(ic.Addr, mux)

	for {
		select {
		case event := <-ic.httpChan:
			inchan <- event
		}
	}
}

// Handler
func (ic *InputConfig) Handler(w http.ResponseWriter, r *http.Request) {
	var (
		message string
		err     error
	)
	r.ParseForm()
	// if r.Method == "GET" {
	for k, v := range r.Form {
		var res []byte
		res, err = json.Marshal(struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{k, strings.Join(v, "")})
		if err != nil {
			glog.Errorln(err)
		}
		message = string(res)
		break
	}
	// }

	event := utils.LogEvent{
		Timestamp: time.Now(),
		Message:   message,
		Extra: map[string]interface{}{
			"host": ic.hostname,
		},
	}
	if err != nil {
		event.AddTag("httpinput error")
	}
	w.Write([]byte(message))
	ic.httpChan <- event
	return
}
