package cmd

import (
	"fmt"
	"logcool/utils"
	_ "logcool/utils/loader" // module loader
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func Logcool(confpath string) (err error) {
	logger := utils.Logger

	if runtime.GOMAXPROCS(0) == 1 && runtime.NumCPU() > 1 {
		logger.Warnf("set GOMAXPROCS = %d to get better performance", runtime.NumCPU())
	}
	fmt.Println("GOMAXPROCS : " + string(runtime.GOMAXPROCS(0)))
	fmt.Println("NumCPU : " + string(runtime.NumCPU()))
	conf, err := utils.LoadFromFile(confpath)
	fmt.Println(confpath)
	fmt.Println(conf)
	if err != nil {
		return
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
		fmt.Println("logcool EXIT...Bye.")
	}
	return
}
