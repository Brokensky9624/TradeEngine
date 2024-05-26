package panichandle

import (
	"fmt"
	"runtime/debug"

	"tradeengine/utils/envmode"
	"tradeengine/utils/logger"
)

func TryCatch(f func()) func() error {
	return func() (err error) {
		defer func() {
			if panicInfo := recover(); panicInfo != nil {
				err = fmt.Errorf("%v, %s", panicInfo, string(debug.Stack()))
				return
			}
		}()
		f() // calling the decorated function
		return err
	}
}

func TryCatchLoop(f func()) func() {
	return func() {
		for {
			if err := TryCatch(f)(); err != nil {
				logger.STD.Error("%v", err)
			} else {
				return
			}
		}
	}
}

func PanicHandle() {
	if panicInfo := recover(); panicInfo != nil {
		logger.Panic.Error("%v, %s", panicInfo, string(debug.Stack()))
		if envmode.IsDevMode() {
			panic(panicInfo)
		}
	}
}
