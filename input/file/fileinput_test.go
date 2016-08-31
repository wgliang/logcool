package fileinput

import (
	"fmt"
	"testing"

	"../../utils"
)

// var (
// 	logger = utils.Logger
// )

// func init() {
// 	logger.Level = logrus.DebugLevel
// 	utils.RegistInputHandler(ModuleName, InitHandler)
// }

// func Test_main(t *testing.T) {
// 	require := require.New(t)
// 	require.NotNil(require)

// 	conf, err := utils.LoadFromString(`{
// 		"input": [{
// 			"type": "file",
// 			"path": "/tmp/log/syslog",
// 			"sincedb_path": "",
// 			"start_position": "beginning"
// 		}]
// 	}`)
// 	require.NoError(err)

// 	err = conf.RunInputs()
// 	require.NoError(err)

// 	waitsec := 10
// 	logger.Infof("Wait for %d seconds", waitsec)
// 	time.Sleep(time.Duration(waitsec) * time.Second)
// }

func Test_DefaultInputConfig(t *testing.T) {
	DefaultInputConfig()
}

func Test_InitHandler(t *testing.T) {
	config := utils.ConfigRaw{}
	co, err := InitHandler(&config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(co)
}

func Test_Start(t *testing.T) {
	// config := DefaultInputConfig()
	// config.Start()
}
