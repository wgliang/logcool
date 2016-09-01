package cmd

import (
	"log"
	"logcool/utils"
	_ "logcool/utils/loader" // module loader
	_ "logcool/utils/logo"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func Logcool(confpath ...string) (err error) {
	logger := utils.Logger
	if runtime.GOMAXPROCS(0) == 1 && runtime.NumCPU() > 1 {
		logger.Warnf("set GOMAXPROCS = %d to get better performance", runtime.NumCPU())
	}
	var conf utils.Config
	if len(confpath) <= 0 {
		conf, err = utils.LoadDefaultConfig()
		if err != nil {
			return
		}
	} else if _, err = os.Stat(confpath[0]); err != nil {
		log.Println("Can not find config-file " + confpath[0] + " and will use default config(stdin2stdout)!")
		conf, err = utils.LoadDefaultConfig()
		if err != nil {
			return
		}
	} else {
		conf, err = utils.LoadFromFile(confpath[0])
		if err != nil {
			log.Println("Config-file " + confpath[0] + " formate error and will use default config(stdin2stdout)!")
			conf, err = utils.LoadDefaultConfig()
			if err != nil {
				return
			}
		}
	}

	if err = conf.RunInputs(); err != nil {
		return
	}

	if err = conf.RunFilters(); err != nil {
		return
	}

	if err = conf.RunOutputs(); err != nil {
		return
	}

	// 捕获ctrl-c,平滑退出
	chExit := make(chan os.Signal, 1)
	signal.Notify(chExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case <-chExit:
		log.Println("logcool EXIT...Bye.")
	}
	return
}
