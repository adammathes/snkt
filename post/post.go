/*
Package post provides the data and behavior for the fundamental atomic
unit of a site: a post. Posts are represented as text files, then converted to HTML and other formats
*/
package post

import (
	"snkt/config"
	"snkt/render"
	"snkt/text"
	"snkt/vlog"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var Template = "post"

type Post struct {
	// Representations of the entire post text
	Raw      []byte
	Unparsed string

	// Metadata
	Meta       map[string]string
	SourceFile string
	Title      string `json:"title"`
	Permalink  string `json:"permalink"`
	Time       time.Time
	Year       int
	Month      time.Month
	Day        int
	InFuture   bool

	// Content text -- raw, unprocessed, unfiltered markdown
	Text string

	// Content text -- processed into HTML via markdown and other filters
	Content string

	// Content with sources and references resolved to absolute URLs
	AbsoluteContent string

	// Post following chronologically (later)
	Next *Post
	// Post preceding chronologically (earlier)
	Prev *Post

	// Precomputed dates as strings
	Date    string
	RssDate string

	FileInfo os.FileInfo

	Site sitemeta
}

type sitemeta interface {
	GetURL() string
	GetTitle() string
}

type Posts []*Post

func (posts Posts) Len() int {
	return len(posts)
}

func (posts Posts) Less(i, j int) bool {
	return posts[i].Time.Before(posts[j].Time)
}

func (s Posts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func NewPost(s sitemeta) *Post {
	var p Post
	p.Site = s
	return &p
}

/*
Read reads a post from file fi, and parses it into the Post struct, performing any work needed to fully populate the struct
*/
func (p *Post) Read(fi os.FileInfo) {
	p.Meta = make(map[string]string)
	p.FileInfo = fi
	p.SourceFile = p.FileInfo.Name()
	var err error
	p.Raw, err = ioutil.ReadFile(path.Join(config.Config.TxtDir, fi.Name()))
	if err != nil {
		log.Println(err)
	}
	p.Unparsed = string(p.Raw)
	p.parse()
}

/*
Parse parses the metadata prefix from the top of the post file's raw bytes, and puts the rest in the text segment. Meta is a name:value mapping
Title, date and other metadata are derived
*/
func (p *Post) parse() {
	//
	// Text + Meta[string][string]
	//
	p.splitTextMeta()

	//
	// Title
	//
	p.Title = p.Meta["title"]
	// Use filename as backup fi we have no explicit title
	if p.Title == "" {
		p.Title = p.SourceFile
	}

	//
	// Dates
	//
	// we only deal with yyyy-mm-dd [some legacy dates from my archives have times tacked on]
	// TODO: recover from empty dates/titles
	// TODO: probably should actually use times when present and clean up my archives
	var date_str = ""
	ds := strings.Fields(p.Meta["date"])
	if len(ds) > 0 {
		date_str = ds[0]
	}

	if date_str == "" {
		p.Time = p.FileInfo.ModTime()
		vlog.Printf("no date field in post %s, using file modification time\n", p.SourceFile)
	} else {
		var err error
		p.Time, err = time.Parse("2006-1-2", date_str)
		if err != nil {
			// fallback is to use file modtime
			// should use create time but that doesn't seem to be in stdlib
			// TODO: figure out how to use file birth time
			vlog.Printf("no valid date parsed for post %s, using file modification time\n", p.SourceFile)
			p.Time = p.FileInfo.ModTime()
		}
	}

	p.Year, p.Month, p.Day = p.Time.Date()
	/* golang date format refresher
	      1 2  3  4  5  7     6
	Mon Jan 2 15:04:05 MST 2006 */

	p.Date = p.Time.Format("January 2, 2006")
	p.RssDate = p.Time.Format("Mon, 2 Jan 2006 15:04:05 GMT")
	p.InFuture = time.Now().Before(p.Time)
	p.Permalink = p.GenPermalink()

	//
	// Content
	//
	p.Content = string(p.Filter([]byte(p.Text)))
	p.AbsoluteContent = render.ResolveURLs(p.Content, p.Site.GetURL())
}

/*
splitText splits up p.Unparsed into p.Text and p.Meta[attr][value]
*/
func (p *Post) splitTextMeta() {
	SEPARATOR := ":"
	lines := strings.Split(p.Unparsed, "\n")
	for _, line := range lines {
		if !strings.Contains(line, SEPARATOR) {
			break
		}
		splitdex := strings.Index(line, SEPARATOR)
		attr := strings.ToLower(strings.TrimSpace(line[0:splitdex]))
		value := strings.TrimSpace(line[splitdex+1:])
		p.Meta[attr] = value
	}
	p.Text = strings.Join(lines[len(p.Meta):], "\n")
}

func (p *Post) ParseFmt(s string) string {
	// TODO: document and add strftime like formats
	s = strings.Replace(s, "%Y", strconv.Itoa(p.Year), -1)
	s = strings.Replace(s, "%M", strconv.Itoa(int(p.Month)), -1)
	s = strings.Replace(s, "%D", strconv.Itoa(p.Day), -1)
	s = strings.Replace(s, "%F", p.CleanFilename(), -1)
	s = strings.Replace(s, "%T", p.CleanTitle(), -1)

	s = strings.Replace(s, "$Y", strconv.Itoa(p.Year), -1)
	s = strings.Replace(s, "$M", strconv.Itoa(int(p.Month)), -1)
	s = strings.Replace(s, "$D", strconv.Itoa(p.Day), -1)
	s = strings.Replace(s, "$F", p.CleanFilename(), -1)
	s = strings.Replace(s, "$T", p.CleanTitle(), -1)

	s = strings.Replace(s, ".File", p.CleanFilename(), -1)
	s = strings.Replace(s, ".Title", p.CleanTitle(), -1)
	s = strings.Replace(s, ".Year", strconv.Itoa(p.Year), -1)
	s = strings.Replace(s, ".Month", strconv.Itoa(int(p.Month)), -1)
	s = strings.Replace(s, ".Day", strconv.Itoa(p.Day), -1)

	return s
}

func (p *Post) CleanFilename() string {
	return text.SanitizeFilename(text.RemoveExt(p.SourceFile))
}

func (p *Post) CleanTitle() string {
	return text.SanitizeFilename(p.Title)
}



/*
GenPermalink generates the permalink for the post given the PermalinkFmt format specified in the configuration file.
*/
func (p *Post) GenPermalink() string {
	pl := config.Config.PermalinkFmt
	return p.ParseFmt(pl)
}

/*
Target returns a string representing the file system location to write the output file representing the post.
*/
func (p Post) Target() string {
	pf := config.Config.PostFileFmt
	return path.Join(config.Config.HtmlDir, p.ParseFmt(pf))
}

/*
Render returns the post rendered as HTML via the post template with Post and Site as context.
*/
func (p Post) Render() []byte {
	data := struct {
		Post interface{}
		Site interface{}
	}{&p, &p.Site}
	return render.Render(Template, data)
}

/*
Filter runs the text through filters defined by render.Filter and marddown, returning text suitable for HTML output.
*/
func (p *Post) Filter(txt []byte) []byte {
	txt = render.Filter(txt)
	txt = blackfriday.MarkdownCommon(txt)
	return txt
}

/*
Limit returns a slice of Posts up to the int limit provided. If the limit is larger than the slice, it just returns the whole slice.
*/
func (posts Posts) Limit(limit int) Posts {
	if len(posts) < limit {
		return posts
	} else {
		return posts[0:limit]
	}
}
