// Output-plug: outputexec
// The plug's function take the event-data into the standard-output.
package outputexec

import (
	"errors"
	"os/exec"

	"github.com/wgliang/logcool/utils"
)

const (
	ModuleName = "lexec"
)

// Define outputexec' config.
type OutputConfig struct {
	utils.OutputConfig
}

func init() {
	utils.RegistOutputHandler(ModuleName, InitHandler)
}

// Init outputexec Handler.
func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeOutputConfig, err error) {
	conf := OutputConfig{
		OutputConfig: utils.OutputConfig{
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

// Input's event,and this is the main function of output.
func (oc *OutputConfig) Event(event utils.LogEvent) (err error) {
	command := event.Message
	if command == "" {
		return errors.New("message is null.")
	}
	args := event.Extra["args"].([]string)
	// run the proc and get all the cmd info.
	Cmd := exec.Command(command, args...)

	// start the commd.
	if err = Cmd.Start(); err != nil {
		return err
	}
	// Wait for the proc done and reset cmd = nil.
	Cmd.Wait()
	return
}
