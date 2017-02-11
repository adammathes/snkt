# snkt

`snkt` is a static web site generator for with a focus on simplicity and efficiency.

snkt only does a few things, but strives to do them well, in a coherent manner.

snkt generates my [personal web site of ~2000 articles in under a second](https://trenchant.org/daily/2017/1/31/). Additional work may be done to increase efficiency, but it should be fast enough to regularly regenerate your site without concern in near real-time if needed.

## Status

Currently in development. It powers [trenchant.org](https://trenchant.org) but is "alpha" quality and parts may change.

TODO:
   * finish these docs
   * lots of other things

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

Install snkt -

    $ go get adammathes.com/snkt

Build `snkt` binary -

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

Build the site with --

    $ snkt -b

Output should now be in the `html` directory and look like -

   * `html`
      * `hi/index.html` hello world post
      * `index.html`
      * `archive.html`
      * `rss.xml` 

You can run a preview server with

    $ snkt -p
    
Loading http://localhost:8000 in a web browser should now show you the site.

Snkt will use `config.yml` by default if it's in the working directory. Otherwise you can specify it with an explicit `-c /path/config.yml` flag.

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

   * `input_dir` -- absolute path of directory for text input files
   * `output_dir` -- absolute path of directory for html output files
   * `tmpl_dir` -- absolute path of directory for template files
   * `site_title` -- string for the site's title (used in templates)
   * `site_url` -- absolute URL for the site (used in templates)
   * `filters` -- tools of the dark arts I haven't documented yet
   * `permalink_fmt` -- format string for permalinks (see #permalinks)
   * `post_file_fmt` -- format string for post filenames (see #permalinks)
   * `show_future` -- include posts with dates in the future or hide them
   * `preview_server` -- host:port to spawn the preview server (default: localhost:8000)
   * `preview_dir` -- root directory of preview server (default: same as output_dir)

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

scary line of unfinished documentation

---

### home
### post
### archive
### rss

## Advanced Features

### paged template

### Permalink and filename formatter

### Filters

### Example configurations/sites/themes

### Auto-rebuild/deployment
