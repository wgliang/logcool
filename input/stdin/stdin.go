package fileinput

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"logcool/utils"
	"logcool/utils/logevent"
)

const (
	ModuleName = "stdin"
)

type InputConfig struct {
	config.InputConfig
}

func DefaultInputConfig() InputConfig {
	return InputConfig{
		InputConfig: config.InputConfig{
			CommonConfig: config.CommonConfig{
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
	fmt.Print("logcool=>")
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "quit" {
			return
		}
		fmt.Println("logcool=>", command)
		fmt.Print("logcool=>")
	}

	event := logevent.LogEvent{
		Timestamp: time.Now(),
		Message:   line,
		Extra: map[string]interface{}{
			"host":   t.hostname,
			"path":   "",
			"offset": since.Offset,
		},
	}

	inchan <- event
	return
}
