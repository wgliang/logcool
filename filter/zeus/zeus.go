//Filter-plug: zues
//zeus is a name of a Greek myth,but it's function is add some fields.Just a easy plug for fun.
package zeus

import (
	"github.com/wgliang/logcool/utils"
	"github.com/wgliang/logcool/utils/logevent"
)

const (
	ModuleName = "zeus"
)

// Define zeus' config.
type FilterConfig struct {
	utils.FilterConfig
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Init zeus Handler.
func InitHandler(confraw *utils.ConfigRaw) (tfc utils.TypeFilterConfig, err error) {
	conf := FilterConfig{
		FilterConfig: utils.FilterConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
	}
	// Reflect config from configraw.
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}

	tfc = &conf
	return
}

// Filter's event,and this is the main function of filter.
func (fc *FilterConfig) Event(event logevent.LogEvent) logevent.LogEvent {
	if _, ok := event.Extra[fc.Key]; ok {
		return event
	}
	if event.Extra == nil {
		event.Extra = make(map[string]interface{})
	}
	event.Extra[fc.Key] = event.Format(fc.Value)
	return event
}
