package log

import (
	"fmt"
	"os"
	"sync/atomic"
)

var verbose int32

func SetVerbose(v bool) {
	if v {
		atomic.StoreInt32(&verbose, 1)
	} else {
		atomic.StoreInt32(&verbose, 0)
	}
}

func Verbose(format string, args ...any) {
	if atomic.LoadInt32(&verbose) == 1 {
		fmt.Fprintf(os.Stderr, "[atlas] "+format+"\n", args...)
	}
}

func Info(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

func Warn(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
}

func Error(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
}
