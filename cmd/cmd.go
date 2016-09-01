// Run Logcool in std, and the filter is zeus.You can input "hello",and it will
// return a formate hello fmt.
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"logcool/utils"
	_ "logcool/utils/loader"
	_ "logcool/utils/logo"
)

// Run Logcool in std, and the filter is zeus.ps:the confpath you can ignore it if you like
func Logcool(confpath ...string) (err error) {
	var conf utils.Config
	// Check the path that you input,if nil will use default config.
	if len(confpath) <= 0 {
		conf, err = utils.LoadDefaultConfig()
		if err != nil {
			return
		}
	} else if _, err = os.Stat(confpath[0]); err != nil {
		fmt.Println("Can not find config-file " + confpath[0] + " and will use default config(stdin2stdout)!")
		conf, err = utils.LoadDefaultConfig()
		if err != nil {
			return
		}
	} else {
		conf, err = utils.LoadFromFile(confpath[0])
		if err != nil {
			fmt.Println("Config-file " + confpath[0] + " formate error and will use default config(stdin2stdout)!")
			conf, err = utils.LoadDefaultConfig()
			if err != nil {
				return
			}
		}
	}
	// Run all Input plugs
	if err = conf.RunInputs(); err != nil {
		return
	}
	// Run all Filter plugs
	if err = conf.RunFilters(); err != nil {
		return
	}
	// Run all Output plugs
	if err = conf.RunOutputs(); err != nil {
		return
	}

	// 捕获ctrl-c,平滑退出
	chExit := make(chan os.Signal, 1)
	signal.Notify(chExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case <-chExit:
		fmt.Println("logcool EXIT...Bye.")
	}
	return
}
