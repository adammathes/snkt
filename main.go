/*
snkt is a static site generator for simple blog-like sites with a focus on simplicity and efficiency.
*/
package main

import (
	"flag"
	"snkt/config"
	"snkt/render"
	"snkt/site"
	"snkt/web"
	"fmt"
	"log"
)

func main() {

	var configFile, init_dir string
	var build, preview, version, verbose, help bool
	
	flag.StringVar(&configFile, "c", "config.yml", "configuration file")
	flag.StringVar(&init_dir, "init", "", "initialize new site at `directory`")
	flag.BoolVar(&build, "b", false, "build the site")
	flag.BoolVar(&preview, "p", false, "start local HTTP server for preview")
	flag.BoolVar(&version, "v", false, "print version number")
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&verbose, "verbose", false, "log more actions while building")
	flag.Parse()

	if !help && !build && !preview && !version && init_dir=="" {
		flag.Usage()
		return
	}
	if(init_dir != "") {
		fmt.Printf("Initializing new site in %s\n", init_dir)
		site.Init(init_dir)
		return
	}
	if(version) {
		fmt.Printf("0.1 alpha\n")
		return
	}
	if(help) {
		fmt.Printf("in case of emergency, break computer \n")
		return
	}
	config.Init(configFile)
	if(verbose) {
		config.Config.Verbose = true
	}

	render.Init()
	if build {
		log.Printf("Building site...\n")
		var s site.Site
		s.Read()
		s.Write()
	}

	if preview {
		log.Printf("Spawning preview at [%s] of [%s]\n",
			config.Config.PreviewDir, config.Config.PreviewServer)
		web.Serve(config.Config.PreviewServer, config.Config.PreviewDir)
	}
}
