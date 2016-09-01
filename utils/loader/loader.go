package loader

import (
	"logcool/filter/zeus"
	"logcool/input/file"
	"logcool/input/stdin"
	"logcool/output/redis"
	"logcool/output/stdout"
	"logcool/utils"
)

func init() {
	utils.RegistInputHandler(fileinput.ModuleName, fileinput.InitHandler)
	utils.RegistInputHandler(stdininput.ModuleName, stdininput.InitHandler)

	utils.RegistFilterHandler(zeus.ModuleName, zeus.InitHandler)

	utils.RegistOutputHandler(outputstdout.ModuleName, outputstdout.InitHandler)
	utils.RegistOutputHandler(outputredis.ModuleName, outputredis.InitHandler)
}
