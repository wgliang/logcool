//Cmd for logcool run.Loading configuration files and information, as well as
//logcool operating logic.
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/wgliang/logcool/utils"
	_ "github.com/wgliang/logcool/utils/logo"
)

const (
	DEFAULTTEMPLATE = "./templates"
	DEFAULTFILE     = "default.json"
)

// load config from Command.
func Command(template string) (confs []utils.Config) {
	conf, err := utils.LoadFromString(template)
	if err != nil {
		log.Println(err)
		return
	}
	confs = append(confs, conf)
	return
}

// load config from user's template.
func Custom(path string) (confs []utils.Config) {
	conf, err := utils.LoadFromFile(path)
	if err != nil {
		log.Println(err)
		return
	}
	confs = append(confs, conf)
	return
}

// load all templates in default dir.
func LoadTemplates() (confs []utils.Config) {
	tempaltes, _ := fileList(DEFAULTTEMPLATE, DEFAULTFILE)
	for _, template := range tempaltes {
		conf, err := utils.LoadFromFile(template)
		if err != nil {
			log.Println(err)
			continue
		} else {
			confs = append(confs, conf)
		}
	}
	if len(confs) <= 0 {
		conf, err := utils.LoadDefaultConfig()
		if err != nil {
			return nil
		} else {
			confs = append(confs, conf)
		}
	}
	return
}

// Run logcool.
func Run(confs []utils.Config) (err error) {
	for _, conf := range confs {
		if err = conf.RunInputs(); err != nil {
			return
		}

		if err = conf.RunFilters(); err != nil {
			return
		}

		if err = conf.RunOutputs(); err != nil {
			return
		}
	}
	return
}

// Help information.
func Help() {
	fmt.Println(`
	-c,-command,  "run in command, stdin2stdout.
	-t,-template, "input templates in command.
	-v,-version,  "show version number.
	-h,-help,     "haha,I know you need me.
	`)
}

// Logcool's version information.
func Version() {
	version, err := ioutil.ReadFile("./VERSION.md")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(version))
}

// list all config templates.
func fileList(dirPath string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dirPath+PthSep+fi.Name())
		}
	}

	return files, nil
}
