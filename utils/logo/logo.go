package logo

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const logo = `
 __                                                ___      
/\ \                                              /\_ \     
\ \ \           ___      __     ___    ___     ___\//\ \    
 \ \ \     __  / __"\  /"_ "\  /"___\ / __"\  / __"\\ \ \   
  \ \ \L___\ \/\ \O\ \/\ \G\ \/\ \C_//\ \O\ \/\ \O\ \\_L \_ 
   \ \_______/\ \____/\ \____ \ \____\ \____/\ \____//\____\
    \/______/  \/___/  \/___L\ \/____/\/___/  \/___/ \/____/
                         /\____/                            
                         \_/__/                             

`
const version = `	Logcool version-0.2`

func info() {
	hostname, err := os.Hostname()
	if err != nil {
		os.Exit(0)
	}

	pid := strconv.Itoa(os.Getpid())
	starttime := time.Now().Format("2006-01-02 03:04:05 PM")
	fmt.Println(version)
	fmt.Println("	Host: " + hostname)
	fmt.Println("	Pid: " + string(pid))
	fmt.Println("	Starttime: " + starttime)
	fmt.Println()
}

func init() {
	fmt.Println(logo)
	info()
}
