// vlog -- verbose logger -- wraps log functions and only performs them if "verbose" config setting is true
package vlog

import (
	"adammathes.com/snkt/config"
	"fmt"
)

func Printf(format string, v ...interface{}) {
	if config.Config.Verbose {
		fmt.Printf(format, v...)
	}
}
	
