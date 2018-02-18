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
	"github.com/fsnotify/fsnotify"
	flag "github.com/ogier/pflag"
)

func main() {

	var configFile, init_dir string
	var build, serve, version, verbose, help, watch bool

	flag.StringVarP(&configFile, "config", "c", "config.yml", "`configuration` file")
	flag.StringVarP(&init_dir, "init", "i", "", "initialize new site at `directory`")
	flag.BoolVarP(&build, "build", "b", false, "generates site from input files and templates")
	flag.BoolVarP(&serve, "serve", "s", false, "serve site via integrated HTTP server")
	flag.BoolVarP(&help, "help", "h", false, "print usage information")
	flag.BoolVarP(&verbose, "verbose", "v", false, "log operations during build to STDOUT")
	flag.BoolVarP(&watch, "watch", "w", false, "watch configured input text dir, rebuild on changes")
	flag.Parse()

	if !watch && !help && !build && !serve && !version && init_dir == "" {
		flag.Usage()
		return
	}
	if init_dir != "" {
		fmt.Printf("initializing new site in %s\n", init_dir)
		site.Init(init_dir)
		return
	}
	if help {
		flag.Usage()
		fmt.Printf("code and docs: https://github.com/adammathes/snkt \n")
		return
	}

	config.Init(configFile)
	if verbose {
		config.Config.Verbose = true
	}

	if build {
		buildSite()
	}

	if serve {
		fmt.Printf("serving [%s] at [%s]\n",
			config.Config.PreviewServer, config.Config.PreviewDir)
		// run on goroutine so watcher can start and block if needed
		go web.Serve(config.Config.PreviewServer, config.Config.PreviewDir)
	}

	if watch {
		watchSite()
	}

	// run preview server perpetually
	if serve {
		done := make(chan bool)
		<-done
	}
}

func buildSite() {
	vlog.Printf("reading templates...\n")
	render.Init()
	var s site.Site
	vlog.Printf("reading posts...\n")
	s.Read()
	vlog.Printf("writing posts and archives...\n")
	s.Write()
	vlog.Printf("done...\n")
}

func watchSite() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("rebuilding... ")
					buildSite()
					fmt.Printf("done\n")
				}
			case err := <-watcher.Errors:
				vlog.Printf("error: %v", err)
			}
		}
	}()

	fmt.Printf("watching %s\n", config.Config.TxtDir)
	err = watcher.Add(config.Config.TxtDir)
	if err != nil {
		panic(err)
	}
	fmt.Printf("watching %s\n", config.Config.TmplDir)
	err = watcher.Add(config.Config.TmplDir)
	if err != nil {
		panic(err)
	}
	<-done
}
