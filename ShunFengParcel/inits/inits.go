package inits

import (
	"os"
	"runtime"
)

func init() {
	InitMysql()

	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
}
