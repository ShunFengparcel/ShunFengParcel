package inits

import (
	"os"
	"runtime"
)

func init() {
	InitMysql()
	InitRedis()

	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
}
