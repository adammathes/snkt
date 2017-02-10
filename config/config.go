package config

import (
	// "encoding/yaml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Settings struct {
	// required
	TxtDir  string `yaml:"input_dir"`
	HtmlDir string `yaml:"output_dir"`
	TmplDir string `yaml:"tmpl_dir"`

	// optional
	SiteTitle string        `yaml:"site_title,omitempty"`
	SiteURL   string        `yaml:"site_url,omitempty"`
	Filters   []RegexFilter `yaml:"filters,omitempty"`

	// required -- set defaults
	PermalinkFmt string `yaml:"permalink_fmt,omitempty"`
	PostFileFmt  string `yaml:"post_file_fmt,omitempty"`
	ShowFuture   bool   `yaml:"show_future,omitempty"`

	PreviewServer string `yaml:"preview_server,omitempty"`
	PreviewDir    string `yaml:"preview_dir,omitempty"`

	Verbose bool `yaml:"verbose,omitempty"`
}

type RegexFilter struct {
	S string `yaml:"s"`
	R string `yaml:"r"`
}

var Config Settings

func Init(filename string) {
	readConfig(filename)
	checkRequired()
	addDefaults()
}

func readConfig(filename string) {
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Fatal("Can not read config file: ", e)
	}
	e = yaml.Unmarshal(file, &Config)
	if e != nil {
		log.Fatal("Config read error: ", e)
	}
}

func checkRequired() {
	if Config.TxtDir == "" {
		log.Fatal("Error: input_dir not set in configuration")
	}
	if Config.HtmlDir == "" {
		log.Fatal("Error: output_dir not set in configuration")
	}
	if Config.TmplDir == "" {
		log.Fatal("Error: tmpl_dir not set in configuration")
	}
}

func addDefaults() {
	if Config.PermalinkFmt == "" {
		Config.PermalinkFmt = "/%F/"
	}
	if Config.PostFileFmt == "" {
		Config.PostFileFmt = "%F/index.html"
	}
	if Config.PreviewServer == "" {
		Config.PreviewServer = "127.0.0.1:8000"
	}
	if Config.PreviewDir == "" {
		Config.PreviewDir = Config.HtmlDir
	}
}
