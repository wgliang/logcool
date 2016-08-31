package main

import (
	"fmt"
	"logcool/cmd"
)

func main() {
	fmt.Println("starting logcool...")
	cmd.Logcool("./templates/stdin2stdout.json")
}
