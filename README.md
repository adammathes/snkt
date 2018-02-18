<pre style="font-family: menlo, courier, monospace;">
                     ██                 
                     ██          ██     
                     ██          ██     
  ▒█████░  ██░████   ██  ▓██▒  ███████  
 ████████  ███████▓  ██ ▓██▒   ███████  
 ██▒  ░▒█  ███  ▒██  ██▒██▒      ██     
 █████▓░   ██    ██  ████▓       ██     
 ░██████▒  ██    ██  █████       ██     
    ░▒▓██  ██    ██  ██░███      ██     
 █▒░  ▒██  ██    ██  ██  ██▒     ██░    
 ████████  ██    ██  ██  ▒██     █████  
 ░▓████▓   ██    ██  ██   ███    ░████  
                                        

              v.01 manual
               2/18/2018

</pre>

- [snkt](#snkt)
    - [Why](#why)
    - [Status](#status)
- [Installation](#installation)
- [Quick Start](#quick-start)
    - [Creating a Site](#creating-a-site)
    - [First Post](#first-post)
    - [Viewing the Results](#viewing-the-results)
- [Command Line Options](#command-line-options)
- [Configuration](#configuration)
- [Posts](#posts)
- [Templates](#templates)
    - [home](#home)
    - [post](#post)
    - [archive](#archive)
    - [rss](#rss)
- [Advanced Configuration Options](#advanced-configuration-options)
    - [Permalink and filename formatter](#permalink-and-filename-formatter)
    - [Filters](#filters)
- [Work in Progress Features](#work-in-progress-features)
    - [Paged Archives](#paged-archives)
    - [Tags](#tags)
    - [Binary Files as Posts](#binary-files-as-posts)
- [Example configurations, sites, themes](#example-configurations-sites-themes)
- [Rebuild and deployment recipes](#rebuild-and-deployment-recipes)
- [TODO](#todo)
- [Feedback](#feedback)

# snkt

`snkt` is a static site generator focused on simplicity and efficiency.

`snkt` does a few things, but strives to do them coherently.

`snkt` generates my [personal web site of ~2000 articles in under a second](https://trenchant.org/daily/2017/1/31/). It should be fast enough to completely regenerate even very large sites in near real-time if needed.

## Why

Every 5-10 years I throw out the software for my site and rewrite it. 

This time it's in Go. Maybe you'll find it useful. It's 10x faster than the old version in Python.

## Status

It powers [trenchant.org](https://trenchant.org) but is under active development and pieces may change. See TODO for future / in progress work.

# Installation

The only dependency for building is `go`.

[Install Go](https://golang.org/doc/install) for your platform.

Download and build `snkt` with something like

    $ go get adammathes.com/snkt

This will download dependencies, build `snkt` and place it in your $GOPATH/bin (by default, ~/go/bin/).

`snkt` is a self-contained binary, you can move it anywhere.

# Quick Start

## Creating a Site

Use the "-init" option to create the skeleton for a new site -

    $ snkt -i myblog
    
This will create:

   * `txt` directory for posts
   * `html` directory for HTML output
   * `tmpl` directory for templates
     * `base` HTML structure wrapper
     * `archive` lists all posts
     * `post` single post page
     * `home` home page with recent posts
     * `rss` RSS 2.0 
   * `config.yml` configuration file

## First Post

A one line plaint text file is a valid post.

    user@host:~/myblog$ echo "hello world" >> txt/hi

Build the site

    $ snkt -b

Output should now be in the `html` directory and look like

   * `html`
      * `hi/index.html` hello world post
      * `index.html`
      * `archive.html`
      * `rss.xml` 

## Viewing the Results

`snkt` includes a simple web server to view the results with

    $ snkt -p
    
Visiting http://localhost:8000 in a web browser should now show the site and the first post.

You can now copy this HTML anywhere and you're set.

# Command Line Options

```
Usage of snkt:
  -b, --build
    	generates site from input files and templates
  -c, --config configuration
    	configuration file (default "config.yml")
  -h, --help
    	print usage information
  -i, --init directory
    	initialize new site at directory
  -s, --serve
    	serve site via integrated HTTP server
  -v, --verbose
    	log operations during build to STDOUT
  -w, --watch
    	watch configured input text dir, rebuild on changes
```

Examples

    $ snkt -c site.yaml -b
    $ snkt --config=myconfig.yml -v -w

# Configuration

Per site configuration is via a [YAML](http://yaml.org) file.

For most purposes, it should just be a listing of attribute : value

Configuration options --

| name             | value                                            | default        |
|------------------|--------------------------------------------------|----------------|
| `input_dir`      | absolute path of directory for text input files  |                |
| `output_dir`     | absolute path of directory for html output files |                |
| `tmpl_dir`       | absolute path of directory for template files    |                |
| `site_title`     | string for the site's title                      |                |
| `site_url`       | absolute URL for the site                        |                |
| `filters`        | list of search/replace regex's to run on posts   |                |
| `permalink_fmt`  | format string for permalinks                     | /%F/           |
| `post_file_fmt`  | format string for post filenames                 | /%F/index.html |
| `show_future`    | include posts with dates in the future           | false          |
| `preview_server` | host:port to spawn the preview server            | localhost:8000 |
| `preview_dir`    | root directory of preview server                 | `output_dir`   |

# Posts

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

# Templates

Templates use the standard library [Go text/template](https://golang.org/pkg/text/template/).

Entities in the templates --

**Site**
```
	Title string
	URL   string
	Posts array of posts
```

See site/site.go for more.

**Post**

```
	// Metadata
	Meta       map[string]string
	SourceFile string
	Title      string 
	Permalink  string
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
```

## home

Displays recent posts and rendered to `index.html` in the `output_dir`.

- {{.Site}} *Site*
- {{.Posts}} *Posts* all posts on site in reverse chronological order

## post

Each individual post uses this template

- {{.Site}} *Site*
- {{.Post}} *Post* the individual post

## archive

Lists all posts, showing only titles and links. Rendered to `archive.html`

- {{.Site}} *Site*
- {{.Posts}} *Posts* all posts, reverse chronological order

## rss

Displays recent posts as RSS 2.0 XML. Rendered to `rss.xml`

- {{.Site}} *Site*
- {{.Posts}} *Posts* all posts, reverse chronological order


# Advanced Configuration Options

## Permalink and filename formatter

Permalinks (URLs for individual posts) can be customized.

| String | Value    | Example |
|--------|----------|---------|
| %Y     | Year     | 2017  |
| %M     | Month    | 04    |
| %D     | Day      | 14    |
| %F     | Filename | foo   |
| %T     | Title    | bar   |

`Filename` is a cleaned version of the post's original filename with the extension removed. Filenames and titles will be "cleaned" of characters unsuitable for links, with whitespace replaced by `-`.


## Filters

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

# Work in Progress Features

These features are working but less documented and potentially still in progress and subject to change.

## Paged Archives

If a template named `paged` is present then paged archives (15 posts per page) are created at `output_dir/page/%d.html`

Template variables are the same as the `archive` template, but with `.NextPage` and `.PrevPage` as integers of the next and previous page.

See archive/paged.go for details.

## Tags

There is preliminary support for tag style metadata per post.

Add a "tags" field to your post preamble. Tags should be comma separated.

```
    tags: TagOne, tag two, a third tag, fourth

```

Tags will be normalized to lowercase, with spaces replaced with underscores. So the above would have tagged a post with --

`tagone tag_two a_third_tag fourth`

Tags are accessible in each post via the `Tags` field.

To create pages by tag, create a template named `tags`.

This creates a file at OUTPUT_DIR/tag/tag_name/index.html for each tag.

It will have access to the same variables as an `archive` template with the additional `.Tag` for the tag name.

## Binary Files as Posts

Preliminary support to treat binary files as standalone posts.

Drop image files with "jpg" or other image extensions into the "txt" dir.

* post's ContentType will be set to "image"
* text fields will be empty strings
* metadata will be populated as it can via exif (maybe)

Video and audio files have preliminary support too -- see `post/post.go`

# Example configurations, sites, themes

*not done*

# Rebuild and deployment recipes

*also not done*

# TODO

   * sample sites/templates
   * proper man pages for docs

# Feedback

Pull requests and issues are welcomed at https://github.com/adammathes/snkt
