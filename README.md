# snkt

`snkt` is a static web site generator for with a focus on simplicity and efficiency.

snkt only does a few things, but strives to do them well, in a coherent manner.

snkt generates my [personal web site of ~2000 articles in under a second](https://trenchant.org/daily/2017/1/31/). Additional work may be done to increase efficiency, but it should be fast enough to regularly regenerate your site without concern in near real-time if needed.

## Status

Currently in development. It powers [trenchant.org](https://trenchant.org) but is "alpha" quality and parts may change.

## TODO

   * finish these docs
   * half-baked / may change
     * permalink formatter
     * filters
   * date handling in templates
   * additional functions in templates
   * themes + example sites
   * complex archive types
   * multiple archives/lists/post outputs

## What

Takes a bunch of plain text files, processes them via templates, and generates HTML. Pretty much like you'd expect of a static site generator.

## Why

Every 5-10 years I throw out the software for my site and rewrite it. 

This time it's in Go. Maybe you'll find it useful. 

I found it fun to get myself thinking in Go. Also, it's 10x faster than the old version in Python.


## Getting snkt

[Install Go](https://golang.org/doc/install)

Set up $GOPATH

    $ mkdir $HOME/go
    $ export GOPATH=$HOME/go
    
Add $GOPATH/bin to your PATH

    $ export PATH=$PATH:$GOPATH/bin

Get and install dependencies

    $ go get gopkg.in/yaml.v2
	$ go get github.com/russross/blackfriday

Download `snkt`

    $ go get adammathes.com/snkt

Build `snkt`

    $ go build adammathes.com/snkt


This should download and build `snkt` and place it in $GOPATH/bin

## Setting up a site

Use the "-init" option to create the skeleton for a new site -

    $ snkt -init blogadu
    
This will create:

   * `txt` directory for plain text input
   * `html` directory for HTML output
   * `tmpl` directory for templates
     * `base` basic HTML structure for all pages
     * `post` single post page template
     * `home` - home page with recent posts template
     * `archive` - list all posts template
     * `rss` - template for an RSS 2.0 archive 
   * `config.yml` -- configuration file


## First Post

A one line plaint text file is a valid post.

    user@host:~/blogadu$ echo "hello world" >> txt/hi

Build the site

    $ snkt -b

Output should now be in the `html` directory and look like

   * `html`
      * `hi/index.html` hello world post
      * `index.html`
      * `archive.html`
      * `rss.xml` 

Run a preview server to see the results with

    $ snkt -p
    
Loading http://localhost:8000 in a web browser should now show the (near empty) site.

## Command Line Usage

```
Usage of snkt:
  -b	build site
  -c configuration
    	configuration file (default "config.yml")
  -h	help
  -init directory
    	initialize new site at directory
  -p	preview site with local HTTP server
  -v	print version number
  -verbose
    	log actions during build
```

## Configuration File

The configuration is in [YAML](http://yaml.org)

For most purposes, it should just be a listing of attribute : value

Configuration options --

| name | value | default |
| `input_dir`      | absolute path of directory for text input files  | |
| `output_dir`     | absolute path of directory for html output files | |
| `tmpl_dir`       | absolute path of directory for template files    | |
| `site_title`     | string for the site's title (used in templates)  | |
| `site_url`       | absolute URL for the site (used in templates)    | |
| `filters`        | search/replace regular expressions executed on all posts | |
| `permalink_fmt`  | format string for permalinks (see #permalinks)           | /%F/ |
| `post_file_fmt`  | format string for post filenames (see #permalinks)  | /%F/index.html |
| `show_future`    | include posts with dates in the future or hide them | false |
| `preview_server` | host:port to spawn the preview server  | localhost:8000 |
| `preview_dir`    | root directory of preview server | `output_dir` |

## Posts

Post inputs are stored as plain text files. (I have only tested UTF-8 and ASCII.)

Posts have an optional metadata preamble, and a markdown formatted body. The preamble is just a series of name value pairings separated by a colon (:) character.

Minimal complete and valid post --

    this is a totally valid post

Post with a preamble --

    title: also a valid post
    date: 2017-02-08
    valid: totes

    This post will have an explicitly set title (ooh! fancy!) 
    instead of inferred from the filename. 

    It will also have an explicitly set date instead of inferring 
    it from the file creation/modification time.

    `totes` will be stored in the post's `meta` map under `valid.` 
    You don't have to worry about that right now. Honest.

## Templates

Templates use the standard library [Go text/template](https://golang.org/pkg/text/template/).

### Types

#### Site (see site/site.go)
```go
	Title string
	URL   string
	Posts post.Posts
```

#### Post (see post/post.go)

```go
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
```

### home

Displays recent posts. Rendered to `index.html`.

- {{.Site}} *Site*
- {{.Posts}} *Posts* all posts on site, reverse chronological order

### post

Template that gets rendered to create individual post pages

- {{.Site}} *Site*
- {{.Post}} *Post* the individual post the page is for

### archive

Lists all posts, showing only titles and links. Rendered to `archive.html`

- {{.Site}} *Site*
- {{.Posts}} *Posts* all posts, reverse chronological order

### rss

Displays recent posts as RSS 2.0 XML. Rendered to `rss.xml`

- {{.Site}} *Site*
- {{.Posts}} *Posts* all posts, reverse chronological order

## Advanced Features

### paged template

If present renders paged archives (15 posts per page) to `page/%d.html`

See archive/paged.go for details. Used to create an "infinite scroll" style archive. Details/options/implementation may change.

### Permalink and filename formatter

Permalink (URLs for individual posts) can be customized. This part is *meh* and subject to change.

| String | Value    | Example |
|--------|----------|---------|
| %Y     | Year     | 2017  |
| %M     | Month    | 04    |
| %D     | Day      | 14    |
| %F     | Filename | foo   |
| %T     | Title    | bar   |

`Filename` is a cleaned version of the post's original filename with the extension removed. Filenames and titles will be "cleaned" of characters unsuitable for links, with whitespace replaced by `-`.

### Filters

Arbitrary regular expressions can be executed on each post to create domain-specific and site-specific modifications.

Here are the real world examples of regular expressions that filter each post on my personal site -

```yaml
filters:
  - s: <photo id="(.+)">
    r: <div class="photo"><img src="/img/$1" /></div>
  - s: <segue />
    r: <p class="segue">&middot; &middot; &middot;</p>
  - s: <youtube id="(.+)">
    r: <p class="video"><a href="https://www.youtube.com/watch?v=$1"><img src="/img/$1.jpg" /></a></p>
  - s: "amazon:(.+)"
    r: "http://www.amazon.com/exec/obidos/ASIN/$1/decommodify-20/"
```

### Example configurations/sites/themes

*not done*

### Auto-rebuild/deployment

*also not done*

