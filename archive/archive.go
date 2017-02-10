package archive

import (
	"adammathes.com/snkt/config"
	"adammathes.com/snkt/post"
	"adammathes.com/snkt/render"
	"path"
)

var archiveTmplName = "archive"
var archiveName = "archive.html"

/*
ListArchive 
*/
type ListArchive struct {
	Posts    post.Posts
	Tgt      string
	Template string
	
	Site interface{}
}

func NewListArchive(posts post.Posts) *ListArchive {
	la := ListArchive{ Posts: posts }
	return &la
}

func (a ListArchive) Target() string {
	if a.Tgt == "" {
		a.Tgt = path.Join(config.Config.HtmlDir, archiveName)
	}
	return a.Tgt
}

func (a ListArchive) Render() []byte {
	if a.Template == "" {
		a.Template = archiveTmplName
	}
	return render.Render(a.Template, a)
}


/*
NewRssArchive takes posts and returns an archive ready for RSS output
*/
func NewRssArchive(posts post.Posts) *ListArchive {
	ra := ListArchive{ Posts: posts, Template: "rss" }
	ra.Tgt = path.Join(config.Config.HtmlDir, "rss.xml")
	return &ra
}
