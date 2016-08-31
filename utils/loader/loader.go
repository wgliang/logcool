package loader

import (
	"fmt"
	"logcool/filter/zeus"
	"logcool/input/file"
	"logcool/output/stdout"
	"logcool/utils"
)

func init() {
	fmt.Println("modulers loader...")
	utils.RegistInputHandler(fileinput.ModuleName, fileinput.InitHandler)

	utils.RegistFilterHandler(zeus.ModuleName, zeus.InitHandler)

	utils.RegistOutputHandler(outputstdout.ModuleName, outputstdout.InitHandler)
}
