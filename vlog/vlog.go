// package vlog wraps log actions and only performs them if "verbose" config is set
package vlog

import (
	"snkt/config"
	"log"
)

func Printf(format string, v ...interface{}) {
	if config.Config.Verbose {
		log.Printf(format, v...)
	}
}
	
