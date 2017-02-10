package site

import (
	"path"
	"adammathes.com/snkt/render"
	"adammathes.com/snkt/config"
	"os"
	"log"
	"gopkg.in/yaml.v2"
)

type skeleton struct {
	Dir string
	Filename string
	Content []byte
}

var skeletons = []skeleton {
	{
			Dir: "tmpl",
			Filename: "base",
			Content: []byte(
`{{define "base"}}
<!DOCTYPE html>
<html>
  <head>
      <title>{{template "title" .}}</title>
  </head>
  <body>
      {{template "content" .}}
      <hr />
      <p><a href="{{.Site.URL}}/index.html">Home</a></p>
      <p><a href="{{.Site.URL}}/archive.html">Archive</a></p>
  </body>
</html>
{{end}}
{{define "title"}}{{end}}
{{define "content"}}{{end}}
`)},
		{
			Dir: "tmpl",
			Filename: "home",
			Content: []byte(
`{{define "title"}}{{.Site.Title}}{{end}}
{{define "content"}}
<h1>{{.Site.Title}}</h1>
{{range .Posts.Limit 15}}
<h2><a href="{{.Permalink}}">{{.Title}}</a></h2>    
{{.Content}}
{{end}}
{{end}}
`)},
	{
			Dir: "tmpl",
			Filename: "post",
			Content: []byte(
`{{define "title"}}{{.Post.Title}} - {{.Site.Title}}{{end}}
{{define "content"}}
<h1>{{.Site.Title}}</h1>
<h2>{{.Post.Title}}</h2>
{{.Post.Content}}
{{if .Post.Next.Title}}
<p><a href="{{.Post.Next.Permalink}}">{{.Post.Next.Title}}</a></p>
{{end}}
    {{if .Post.Prev.Title}}
<p><a href="{{.Post.Prev.Permalink}}">{{.Post.Prev.Title}}</a></p>
{{end}}
{{end}}
`)},
		{
			Dir: "tmpl",
			Filename: "archive",
			Content: []byte(
`{{define "title"}}{{.Site.Title}} Archives{{end}}
{{define "content"}}
<h1>{{ .Site.Title }}</h1>
<h2>Archives</h2>
{{range .Posts}}
{{if .}}<h2><a href="{{.Permalink}}">{{.Title}}</a></h2>{{end}}
{{end}}
{{end}}
`)},
		{
			Dir: "tmpl",
			Filename: "rss",
			Content: []byte(
`{{define "base"}}
<?xml version="1.0" encoding="UTF-8" ?> 
<rss version="2.0">
<channel>
  <title>{{ .Site.Title }}</title>
  <link>{{ .Site.URL }}</link> 
  <description></description>
  {{ range .Posts.Limit 15 }}
  <item>    
    <link>{{ .Permalink }}</link>
    <title><![CDATA[{{ .Title }}]]></title>
    <pubDate>{{ .RssDate }}</pubDate>
    <description><![CDATA[ 
    {{ .AbsoluteContent }}
    ]]>
    </description>
  </item>
  {{ end }}
</channel>
</rss>
{{end}}
`)}}

/*
Init initializes a new site in `directory` This includes: create `directory` and populate with:
config.json with sane defaults
txt, html directories for input/output
tmpl directory with barebones [base,post,archive,rss] templates
*/
var init_dir = ""
func Init(directory string) {
	init_dir = directory

	if init_dir[0] != '/' {
		wd, err := os.Getwd()
		if err == nil {
			init_dir = path.Join(wd, init_dir)
		}
	}

	var cfg = config.Settings{
		SiteTitle: "snkt",
		SiteURL: "",
		TxtDir: path.Join(init_dir, "txt"),
		HtmlDir: path.Join(init_dir, "html"),
		TmplDir: path.Join(init_dir, "tmpl"),
	}

	cyaml, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatal("marshalling yaml error: ", err)
	}
	
	c := skeleton{
		Dir: "",
		Filename: "config.yml",
		Content: cyaml,
	}
	skeletons = append(skeletons, c)

	os.MkdirAll( cfg.TxtDir , 0755)
	os.MkdirAll( cfg.HtmlDir , 0755)
	os.MkdirAll( cfg.TmplDir , 0755)

	writeSkeletons()
}


func (s skeleton) Render() []byte {
	return s.Content
}

func (s skeleton) Target() string {
	return path.Join(init_dir, s.Dir, s.Filename)
}

func writeSkeletons() {
	for _,skeleton := range skeletons {
		render.Write(skeleton)
	}
}
