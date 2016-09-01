package main

import (
	"flag"
	"fmt"
	"github.com/wgliang/logcool/cmd"
	"github.com/wgliang/logcool/utils"
	"github.com/wgliang/logcool/utils/config"
	"os"
	"os/signal"
	"syscall"
)

var (
	conf    = flag.String("config", "", "path to config.json file")
	command = flag.String("command", "", "run in command, stdin2stdout.")
	custom  = flag.String("custom", "", "input custom template.")
	version = flag.Bool("version", false, "show version number.")
	std     = flag.Bool("std", false, "run in stadin/stdout.")
	help    = flag.Bool("help", false, "haha,I know you need me.")
)

func main() {
	flag.Parse()

	if *version != false {
		config.Version()
		os.Exit(0)
	}

	if *help != false {
		config.Help()
		os.Exit(0)
	}

	if *std != false {
		cmd.Logcool()
		os.Exit(0)
	}

	var confs []utils.Config
	if *custom != "" {
		confs = config.Custom(*custom)
	} else if *command != "" {
		confs = config.Command(*command)
	} else {
		confs = config.LoadTemplates()
	}

	config.Run(confs)

	// 捕获ctrl-c,平滑退出
	chExit := make(chan os.Signal, 1)
	signal.Notify(chExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case <-chExit:
		fmt.Println("logcool EXIT...Bye.")
	}
}
