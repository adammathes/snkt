package site

import (
	"adammathes.com/snkt/archive"
	"adammathes.com/snkt/config"
	"adammathes.com/snkt/post"
	"adammathes.com/snkt/render"
	"io/ioutil"
	"log"
	"path"
	"sort"
)

type Site struct {
	Title string
	URL   string

	Posts post.Posts

	// all archives are optional based on presence of template
	Archive *archive.ListArchive
	Home    *archive.ListArchive
	Rss     *archive.ListArchive
	Paged   *archive.PagedArchives
}

/*
Read reads post data from the filesystem and populates posts and archives
*/
func (s *Site) Read() {
	s.Title = config.Config.SiteTitle
	s.URL = config.Config.SiteURL
	s.ReadPosts()

	if render.TmplExists("archive") {
		s.Archive = archive.NewListArchive(s.Posts)
		s.Archive.Site = *s
		sort.Sort(sort.Reverse(s.Archive.Posts))
	}
	if render.TmplExists("rss") {
		s.Rss = archive.NewRssArchive(s.Posts)
		s.Rss.Site = *s
	}
	if render.TmplExists("paged") {
		s.Paged = archive.CreatePaged(15, s.Posts)
	}
	if render.TmplExists("home") {
		s.Home = archive.NewListArchive(s.Posts)
		s.Home.Tgt = path.Join(config.Config.HtmlDir, "index.html")
		s.Home.Template = "home"
		s.Home.Site = s
	}
}

/*
ReadPosts reads all files from the Config.TxtDir, parses them and stores in s.Posts 
*/
func (s *Site) ReadPosts() {
	// TODO: filter this as needed
	files, err := ioutil.ReadDir(config.Config.TxtDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		p := post.NewPost(s)
		p.Read(file)

		// Ignore post-dated posts unless overriden in config
		if !config.Config.ShowFuture && p.InFuture {
			log.Printf("Skipping future dated post: %s\n", p.SourceFile)
		} else {
			s.Posts = append(s.Posts, p)
		}
	}

	// Sort the posts by date, earliest first
	sort.Sort(s.Posts)

	// Add next/previous to each post.
	// An allocated but empty post is set at start/end
	// This prevents templates from failing at start/end if nil checks are not made properly
	for i, p := range s.Posts {
		if i > 0 {
			p.Prev = s.Posts[i-1]
		} else {
			p.Prev = new(post.Post)
		}
		if i+1 < len(s.Posts) {
			p.Next = s.Posts[i+1]
		} else {
			p.Next = new(post.Post)
		}
	}
}

/*
Write writes posts and archives to the filesystem
*/
func (s *Site) Write() {
	s.WritePosts()
	s.WriteArchives()
}

func (s *Site) WriteArchives() {
	if render.TmplExists("archive") {
		render.Write(s.Archive)
	}
	if render.TmplExists("rss") {
		render.Write(s.Rss)
	}
	if render.TmplExists("home") {
		render.Write(s.Home)
	}
	if render.TmplExists("paged") {
		for _, p := range *s.Paged {
			p.Site = s
			render.Write(p)
		}
	}
}

func (s *Site) WritePosts() {
	for _, p := range s.Posts {
		render.Write(p)
	}
}

func (s Site) GetTitle() string {
	return s.Title
}

func (s Site) GetURL() string {
	return s.URL
}
