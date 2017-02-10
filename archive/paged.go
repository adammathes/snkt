package archive

import (
	"adammathes.com/snkt/config"
	"adammathes.com/snkt/post"
	"adammathes.com/snkt/render"
	"fmt"
	"math"
	"path"
	"sort"
)

var pagedTmplName = "paged"

/*
Paged archive shows set of posts broken up over multiple pages
Output goes to Config.HtmlDir/page/{pageNum}.html
*/
type PagedArchive struct {
	Posts    post.Posts
	PageNum  int
	NextPage int
	PrevPage int

	Site interface{}
}

func (pa PagedArchive) Render() []byte {
	return render.Render(pagedTmplName, pa)
}

/* TODO: make this configurable */
func (pa PagedArchive) Target() string {
	return path.Join(config.Config.HtmlDir, "/page/", fmt.Sprintf("%d.html", pa.PageNum))
}

type PagedArchives []*PagedArchive

func CreatePaged(perPage int, posts post.Posts) *PagedArchives {

	if !render.TmplExists(pagedTmplName) {
		fmt.Printf("no page template\n")
		return nil
	}
	
	var pas PagedArchives

	sort.Sort(sort.Reverse(posts))

	numPages := int(math.Ceil(float64(len(posts)) / float64(perPage)))
	for i := 0; i < numPages; i++ {
		var pa PagedArchive

		var m int
		if (i+1)*perPage > len(posts) {
			m = len(posts)
		} else {
			m = (i + 1) * perPage
		}

		pa.Posts = posts[i*perPage : m]
		pa.PrevPage = i
		pa.PageNum = i + 1
		pa.NextPage = i + 2
		pas = append(pas, &pa)
	}
	return &pas
}

func (pas *PagedArchives) Write() {
	for _, pa := range *pas {
		render.Write(pa)
	}
}
