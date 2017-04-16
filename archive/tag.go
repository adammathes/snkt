package archive

import (
	"adammathes.com/snkt/config"
	"adammathes.com/snkt/post"
	"adammathes.com/snkt/render"
	"fmt"
	"path"
)

var tagTmplName = "tag"

/*
Tag archive shows set of posts broken by tag
Output goes to Config.HtmlDir/tag/{tag}.html
*/
type TagArchive struct {
	Posts post.Posts
	Tag   string
	Site  interface{}
}

func (ta TagArchive) Render() []byte {
	return render.Render(tagTmplName, ta)
}

func (ta TagArchive) Target() string {
	return path.Join(config.Config.HtmlDir, "tag", ta.Tag, "index.html")
}

type TagArchives []*TagArchive

func ParseTags(posts post.Posts) *TagArchives {

	if !render.TmplExists(tagTmplName) {
		fmt.Printf("no tag template\n")
		return nil
	}

	var tas TagArchives

	// create a map of [tag]posts
	var tags map[string]post.Posts
	tags = make(map[string]post.Posts)

	for _, p := range posts {
		for _, t := range p.Tags {
			_, ok := tags[t]
			if !ok {
				var ps post.Posts
				tags[t] = ps
			}
			tags[t] = append(tags[t], p)
		}
	}

	for tag, posts := range tags {
		var ta TagArchive
		ta.Tag = tag
		ta.Posts = posts
		tas = append(tas, &ta)
	}

	return &tas
}

func (tas *TagArchives) Write() {
	for _, ta := range *tas {
		render.Write(ta)
	}
}
