package render

import (
	"bytes"
	"adammathes.com/snkt/config"
	"io/ioutil"
	"log"
	"adammathes.com/snkt/vlog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"text/template"
)

var templates map[string]*template.Template
var BASE_TEMPLATE = "base"
var rel_href *regexp.Regexp
var rel_src *regexp.Regexp

/*
Renderable interface - objects that render themeslves to a []byte and
know where they should end up in the filesystem
*/
type Renderable interface {
	Render() []byte
	Target() string
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
		"SiteTitle": SiteTitle,
		"SiteURL": SiteURL,
	}

	
	base := path.Join(config.Config.TmplDir, BASE_TEMPLATE)
	for _, t := range ts {
		tf := filepath.Base(t)

		// Duly Noted: set funcs before parsefiles or you get no funcs
		// templates[tf] = template.Must(template.ParseFiles(t, base))

		tx := template.New("t").Funcs(tmplFuncs)
		templates[tf], err = tx.ParseFiles(base, t)
		if err != nil {
			panic(err)
		}
	}
	rel_href = regexp.MustCompile(`href="/(.+)"`)
	rel_src = regexp.MustCompile(`src="/(.+)"`)
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
