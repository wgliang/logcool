package zeus

import (
	"logcool/utils"
	"logcool/utils/logevent"
)

const (
	ModuleName = "zeus"
)

type FilterConfig struct {
	utils.FilterConfig
	Key   string `json:"key"`
	Value string `json:"value"`
}

func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeFilterConfig, err error) {
	conf := FilterConfig{
		FilterConfig: utils.FilterConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
	}
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}

	retconf = &conf
	return
}

func (f *FilterConfig) Event(event logevent.LogEvent) logevent.LogEvent {
	if _, ok := event.Extra[f.Key]; ok {
		return event
	}
	if event.Extra == nil {
		event.Extra = make(map[string]interface{})
	}
	event.Extra[f.Key] = event.Format(f.Value)
	return event
}
