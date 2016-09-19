//Filter-plug: split
//split is just for split and count.
package split

import (
	"strings"

	"github.com/wgliang/logcool/utils"
)

const (
	ModuleName = "split"
)

// Define split' config.
type FilterConfig struct {
	utils.FilterConfig
	Separator string `json:"separator"`
}

func init() {
	utils.RegistFilterHandler(ModuleName, InitHandler)
}

// Init split Handler.
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
func (fc *FilterConfig) Event(event utils.LogEvent) utils.LogEvent {
	if event.Extra == nil {
		event.Extra = make(map[string]interface{})
	}
	args := strings.Split(event.Message, fc.Separator)
	if len(args) > 0 {
		event.Message = args[0]
	}

	event.Extra["args"] = args[1:]

	return event
}
