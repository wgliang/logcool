package outputstdout

import (
	"fmt"

	"logcool/utils"
	"logcool/utils/logevent"
)

const (
	ModuleName = "stdout"
)

type OutputConfig struct {
	utils.OutputConfig
}

func DefaultOutputConfig() OutputConfig {
	return OutputConfig{
		OutputConfig: utils.OutputConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
	}
}

func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeOutputConfig, err error) {
	conf := DefaultOutputConfig()
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}

	retconf = &conf
	return
}

func (t *OutputConfig) Event(event logevent.LogEvent) (err error) {
	raw, err := event.MarshalIndent()
	if err != nil {
		return
	}

	fmt.Println(string(raw))
	return
}
