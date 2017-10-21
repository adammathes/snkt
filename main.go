/*
snkt is a static site generator for simple blog-like sites with a focus on simplicity and efficiency.
*/
package main

import (
	"adammathes.com/snkt/config"
	"adammathes.com/snkt/render"
	"adammathes.com/snkt/site"
	"adammathes.com/snkt/vlog"
	"adammathes.com/snkt/web"
	"fmt"
	flag "github.com/ogier/pflag"
)

func main() {

	var configFile, init_dir string
	var build, preview, version, verbose, help bool

	flag.StringVarP(&configFile, "config", "c", "config.yml", "`configuration` file")
	flag.StringVarP(&init_dir, "init", "i", "", "initialize new site at `directory`")
	flag.BoolVarP(&build, "build", "b", false, "build site")
	flag.BoolVarP(&preview, "preview", "p", false, "preview site with local HTTP server")
	flag.BoolVarP(&help, "help", "h", false, "print help message")
	flag.BoolVarP(&verbose, "verbose", "v", false, "log actions during build to STDOUT")
	flag.Parse()

	if !help && !build && !preview && !version && init_dir == "" {
		flag.Usage()
		return
	}
	if init_dir != "" {
		fmt.Printf("initializing new site in %s\n", init_dir)
		site.Init(init_dir)
		return
	}
	if help {
		fmt.Printf("please see README.md\n")
		return
	}
	config.Init(configFile)
	if verbose {
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
