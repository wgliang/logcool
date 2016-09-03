package cmd

import (
	"os"
	"testing"
	"time"
)

func Test_Logcool(t *testing.T) {
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
	Logcool()
}
