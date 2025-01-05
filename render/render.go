package render

import (
	"adammathes.com/snkt/config"
	"adammathes.com/snkt/vlog"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"text/template"
)

var templates map[string]*template.Template
var BASE_TEMPLATE = "base"
var rel_href, rel_src, re_href, re_src *regexp.Regexp

/*
Renderable interface - objects that render themeslves to a []byte and
know where they should end up in the filesystem
*/
type Renderable interface {
	Render() []byte
	Target() string
}

func TemplateNames() []string {
	templateNames := make([]string, len(templates))

	i := 0
	for tName, _ := range templates {
		templateNames[i] = tName
		i++
	}

	return templateNames
}

func Write(a Renderable) {
	if config.Config.Verbose {
		vlog.Printf("Writing to %s\n", a.Target())
	}
	os.MkdirAll(path.Dir(a.Target()), 0755)
	err := ioutil.WriteFile(a.Target(), a.Render(), 0755)
	if err != nil {
		log.Println(err)
	}
}

/*
Initializes templates from config.TmplDir
Templates are mapped by filename
*/
func Init() {
	templates = make(map[string]*template.Template)
	ts, err := filepath.Glob(config.Config.TmplDir + "/*")
	if err != nil {
		log.Fatal(err)
	}

	tmplFuncs := template.FuncMap{
		"ResolveURLs": ResolveURLs,
		"SiteTitle":   SiteTitle,
		"SiteURL":     SiteURL,
	}

	// base := path.Join(config.Config.TmplDir, BASE_TEMPLATE)
	for _, t := range ts {
		tf := filepath.Base(t)

		// Duly Noted: set funcs before parsefiles or you get no funcs
		// templates[tf] = template.Must(template.ParseFiles(t, base))

		tx := template.New("t").Funcs(tmplFuncs)
		//		templates[tf], err = tx.ParseFiles(base, t, (ts...))
		templates[tf], err = tx.ParseGlob(config.Config.TmplDir + "/*")
		if err != nil {
			// temporary files can confuse this, especially when
			// running the file system watcher so we silently
			// ignore any templates that disappeared since we started
			// since this is usually not a real  error condition
			delete(templates, tf)
		}
	}
	rel_href = regexp.MustCompile(`href="/(.+)"`)
	rel_src = regexp.MustCompile(`src="/(.+)"`)
	re_href = regexp.MustCompile(`href="(.*?)"`)
	re_src = regexp.MustCompile(`src="(.*?)"`)
}

/*
Render fills the template "name" using data via the BASE_TEMPLATE
*/
func Render(name string, data interface{}) []byte {
	return RenderNameVia(name, BASE_TEMPLATE, data)
}

/*
Render fills the template "name" using data
*/
func RenderOnly(name string, data interface{}) []byte {
	return RenderNameVia(name, name, data)
}

/*
Render fills the template "name" using data through via (ex: BASE_TEMPLATE)
*/
func RenderNameVia(name string, via string, data interface{}) []byte {
	t, ok := templates[name]
	if !ok {
		log.Printf("can not find template named %s\n", name)
	}

	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf, via, data)
	if err != nil {
		log.Println(err)
	}
	return buf.Bytes()
}

/*
Finds any relative links/images in html and resolves by adding prefix
*/
func ResolveURLs(html, prefix string) string {
	bts := []byte(html)
	bts = rel_href.ReplaceAll(bts, []byte(`href="`+prefix+`/$1"`))
	bts = rel_src.ReplaceAll(bts, []byte(`src="`+prefix+`/$1"`))
	return string(bts)
}

/*
Finds all URLs that are hrefs
TODO: replace noisy regex with HTML parser
*/
func FindURLs(html string) []string {
	//	bts := []byte(html)
	hrefs := re_href.FindAllStringSubmatch(html, -1)
	var urls []string
	for _, href := range hrefs {
		urls = append(urls, href[1])
	}
	return urls
}

/*
Finds all img urls via img src tags
TODO: replace noisy regex with HTML parser
*/
func FindImgs(html string) []string {
	//	bts := []byte(html)
	srcs := re_src.FindAllStringSubmatch(html, -1)
	var imgs []string
	for _, src := range srcs {
		imgs = append(imgs, src[1])
	}
	return imgs
}

/*
Runs all regex filters specified in config.Config.Filters
*/
func Filter(txt []byte) []byte {
	for _, f := range config.Config.Filters {
		// TODO: only compile these once at init
		re := regexp.MustCompile(f.S)
		txt = re.ReplaceAll(txt, []byte(f.R))
	}
	return txt
}

func TmplExists(t string) bool {
	_, ok := templates[t]
	return ok
}

func SiteTitle() string {
	return config.Config.SiteTitle
}

func SiteURL() string {
	return config.Config.SiteURL
}
