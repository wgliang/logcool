package stdininput

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"logcool/utils"
	"logcool/utils/logevent"
)

const (
	ModuleName = "stdin"
)

type InputConfig struct {
	utils.InputConfig

	hostname string `json:"-"`
}

func DefaultInputConfig() InputConfig {
	return InputConfig{
		InputConfig: utils.InputConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
	}
}

func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeInputConfig, err error) {
	conf := DefaultInputConfig()
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}

	if conf.hostname, err = os.Hostname(); err != nil {
		return
	}

	retconf = &conf
	return
}

func (t *InputConfig) Start() {
	t.Invoke(t.start)
}

func (t *InputConfig) start(logger *logrus.Logger, inchan utils.InChan) (err error) {
	defer func() {
		if err != nil {
			logger.Errorln(err)
		}
	}()

	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		time.Sleep(300000 * time.Nanosecond)
		fmt.Print("logcool#")
		data, _, _ := reader.ReadLine()
		command := string(data)
		event := logevent.LogEvent{
			Timestamp: time.Now(),
			Message:   command,
			Extra: map[string]interface{}{
				"host": t.hostname,
			},
		}
		inchan <- event
		if command == "quit" {
			os.Exit(0)
		}
	}

	return
}
