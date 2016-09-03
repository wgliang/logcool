package loader

import (
	"github.com/wgliang/logcool/filter/zeus"
	"github.com/wgliang/logcool/input/file"
	"github.com/wgliang/logcool/input/http"
	"github.com/wgliang/logcool/input/stdin"
	"github.com/wgliang/logcool/output/redis"
	"github.com/wgliang/logcool/output/stdout"
	"github.com/wgliang/logcool/utils"
)

func init() {
	utils.RegistInputHandler(fileinput.ModuleName, fileinput.InitHandler)
	utils.RegistInputHandler(stdininput.ModuleName, stdininput.InitHandler)
	utils.RegistInputHandler(httpinput.ModuleName, httpinput.InitHandler)

	utils.RegistFilterHandler(zeus.ModuleName, zeus.InitHandler)

	utils.RegistOutputHandler(outputstdout.ModuleName, outputstdout.InitHandler)
	utils.RegistOutputHandler(outputredis.ModuleName, outputredis.InitHandler)
}
