package outputstdout

import (
	"reflect"
	"testing"
	"time"

	"../../utils"
	"../../utils/logevent"
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	logger = utils.Logger
)

func init() {
	logger.Level = logrus.DebugLevel
	utils.RegistOutputHandler(ModuleName, InitHandler)
}

func Test_main(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	conf, err := utils.LoadFromString(`{
		"output": [{
			"type": "stdout"
		}]
	}`)
	require.NoError(err)

	err = conf.RunOutputs()
	require.NoError(err)

	outchan := conf.Get(reflect.TypeOf(make(utils.OutChan))).
		Interface().(utils.OutChan)
	outchan <- logevent.LogEvent{
		Timestamp: time.Now(),
		Message:   "outputstdout test message",
	}

	waitsec := 1
	logger.Infof("Wait for %d seconds", waitsec)
	time.Sleep(time.Duration(waitsec) * time.Second)
}

func Test_DefaultOutputConfig(t *testing.T) {
	DefaultOutputConfig()
}

func Test_InitHandler(t *testing.T) {
	// InitHandler()
}

func Test_Event(t *testing.T) {
	le := logevent.LogEvent{
		Timestamp: time.Now(),
		Message:   "message",
		Tags:      []string{"frg", "grbhrt"},
		Extra:     make(map[string]interface{}),
	}
	config := DefaultOutputConfig()
	config.Event(le)
}
