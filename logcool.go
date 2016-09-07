package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wgliang/logcool/cmd"
	_ "github.com/wgliang/logcool/filter/zeus"
	_ "github.com/wgliang/logcool/input/collectd"
	_ "github.com/wgliang/logcool/input/file"
	_ "github.com/wgliang/logcool/input/http"
	_ "github.com/wgliang/logcool/input/stdin"
	_ "github.com/wgliang/logcool/output/redis"
	_ "github.com/wgliang/logcool/output/stdout"
	"github.com/wgliang/logcool/utils"
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
		cmd.Version()
		os.Exit(0)
	}

	if *help != false {
		cmd.Help()
		os.Exit(0)
	}
	var confs []utils.Config

	if *std != false {
		// cmd.Logcool()
		conf, err := utils.LoadDefaultConfig()
		if err != nil {
			fmt.Println(err)
		}
		confs = append(confs, conf)
	} else if *custom != "" {
		confs = cmd.Custom(*custom)
	} else if *command != "" {
		confs = cmd.Command(*command)
	} else {
		confs = cmd.LoadTemplates()
	}

	cmd.Run(confs)

	// 捕获ctrl-c,平滑退出
	chExit := make(chan os.Signal, 1)
	signal.Notify(chExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case <-chExit:
		fmt.Println("logcool EXIT...Bye.")
	}
}
