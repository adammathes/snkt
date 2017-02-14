/*
snkt is a static site generator for simple blog-like sites with a focus on simplicity and efficiency.
*/
package main

import (
	"flag"
	"adammathes.com/snkt/config"
	"adammathes.com/snkt/render"
	"adammathes.com/snkt/site"
	"adammathes.com/snkt/web"
	"adammathes.com/snkt/vlog"
	"fmt"
)

func main() {

	var configFile, init_dir string
	var build, preview, version, verbose, help bool
	
	flag.StringVar(&configFile, "c", "config.yml", "`configuration` file")
	flag.StringVar(&init_dir, "init", "", "initialize new site at `directory`")
	flag.BoolVar(&build, "b", false, "build site")
	flag.BoolVar(&preview, "p", false, "preview site with local HTTP server")
	flag.BoolVar(&version, "v", false, "print version number")
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&verbose, "verbose", false, "log actions during build")
	flag.Parse()

	if !help && !build && !preview && !version && init_dir=="" {
		flag.Usage()
		return
	}
	if(init_dir != "") {
		fmt.Printf("initializing new site in %s\n", init_dir)
		site.Init(init_dir)
		return
	}
	if(version) {
		fmt.Printf("0.1 alpha\n")
		return
	}
	if(help) {
		fmt.Printf("please see README.md\n")
		return
	}
	config.Init(configFile)
	if(verbose) {
		config.Config.Verbose = true
	}

	render.Init()
	if build {
		var s site.Site
		vlog.Printf("reading posts...\n")
		s.Read()
		vlog.Printf("writing posts and archives...\n")
		s.Write()
	}

	if preview {
		fmt.Printf("spawning preview at [%s] of [%s]\n",
			config.Config.PreviewDir, config.Config.PreviewServer)
		web.Serve(config.Config.PreviewServer, config.Config.PreviewDir)
	}
}
