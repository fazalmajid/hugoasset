package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"path/filepath"
	"fmt"
	"strings"
	"time"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"github.com/lukasbob/srcset"
	//"github.com/jaytaylor/html2text"
	//"github.com/gohugoio/hugo/parser"
	//"github.com/gohugoio/hugo/hugolib"
	//"github.com/spf13/cast"
	//"github.com/spf13/afero"
	//"github.com/gohugoio/hugo/hugofs"
	//"github.com/gohugoio/hugo/deps"
)

const date_fmt = "[02/Jan/2006:15:04:05 -0700]"
const iso_8601 = "2006-01-02 15:04:05"

var (
	verbose *bool
	prefix *string
)

func main() {
	// command-line options
	verbose = flag.Bool("v", false, "Verbose error reporting")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	do_html := flag.Bool("html", true, "Search HTML files for assets")
	prefix = flag.String("prefix", "/blog", "site prefix for assets")
	//do_hugo := flag.Bool("hugo", false, "Search Hugo markdown files for assets")
	do_unused := flag.Bool("unused", true, "Search for unused assets")
	flag.Parse()
	var err error
	var f *os.File
	// Profiler
	if *cpuprofile != "" {
		f, err = os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var updated time.Time

	if *do_html {
		before := time.Now()
		log.Println("indexing HTML...")
		seen := index_html(updated)
		if *do_unused {
			find_unused(seen)
		}
		log.Println("done in", time.Now().Sub(before))
	}
	// if *do_hugo {
	// 	before := time.Now()
	// 	log.Println("indexing Hugo...")
	// 	index_hugo(db, updated)
	// 	log.Println("done in", time.Now().Sub(before))
	// }
}

// recurse through the parsed HTML tree to extract the title
func extract_urls(n *html.Node) []string {
	urls := make([]string, 0)
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		for _, a := range(n.Attr) {
			if a.Key == "href" {
				urls = append(urls, a.Val)
				break
			}
		}
		return urls
	}
	if n.Type == html.ElementNode && n.DataAtom == atom.Img {
		for _, a := range(n.Attr) {
			if a.Key == "src" {
				urls = append(urls, a.Val)
			}
			if a.Key == "srcset" {
				ss := srcset.Parse(a.Val)
				if ss != nil {
					for _, img := range(ss) {
						urls = append(urls, img.URL)
					}
				}
			}
		}
		return urls
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := extract_urls(c)
		if result != nil {
			urls = append(urls, result...)
		}
	}
	return urls
}

func index_html(updated time.Time) map[string]bool {
	seen := make(map[string]bool, 100)
	// walk the current directory looking for HTML files
	err := filepath.Walk("public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fn := strings.ToLower(path)
		if !(strings.HasSuffix(fn, ".html") || strings.HasSuffix(fn, ".htm")) {
			return nil
		}
		if strings.Contains(path, "/categories/") {
			return nil
		}
		if strings.Contains(path, "/tags/") {
			return nil
		}
		stat, err := os.Stat(path)
		if stat.ModTime().Before(updated) {
			return nil
		}
		if *verbose {
			fmt.Println(path);
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		html_doc, err := html.Parse(f)
		if err != nil {
			return err
		}
		f.Close()
		urls := extract_urls(html_doc)
		printed := false
		for _, u := range(urls) {
			if u[len(u)-1] == '/' {
				u = u + "index.html"
			}
			if strings.HasPrefix(u, *prefix) {
				unescaped, err := url.QueryUnescape(u[len(*prefix):len(u)])
				if err == nil {
					_, err = os.Stat("public/" + unescaped)
				}
				if err != nil {
					if !printed {
						fmt.Println(path)
						printed = true
					}
					fmt.Println("\t", u)
					fmt.Println("\t\t", err)
				} else {
					seen[unescaped] = true
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error while walking html", err)
	}
	return seen
}

func find_unused(seen map[string]bool) error {
	if *verbose {
		log.Println("looking for unused assets")
	}
	// walk the current directory looking for HTML files
	err := filepath.Walk("static", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ! info.Mode().IsRegular() {
			return nil
		}
		if *verbose {
			log.Println("\t", path)
		}
		if len(path) < len("static") {
			return nil
		}
		fn := path[len("static"):len(path)]
		_, ok := seen[fn]
		if !ok {
			fmt.Println("unused asset:\t", path)
		}
		return nil
	})
	return err
}

// func index_hugo(db *sql.DB, updated time.Time) {
// 	osFs := hugofs.Os
// 	cfg, err := hugolib.LoadConfig(osFs, "", "config.toml")
// 	if err != nil {
// 		log.Fatal("Could not load Hugo config.toml", err)
// 	}
// 	fs := hugofs.NewFrom(osFs, cfg)
// 	sites, err := hugolib.NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})
// 	if err != nil {
// 		log.Fatal("Could not load Hugo site(s)", err)
// 	}
// 	err =sites.Build(hugolib.BuildCfg{SkipRender: true})
// 	if err != nil {
// 		log.Fatal("Could not run render", err)
// 	}
// 	for _, p := range sites.Pages() {
// 		if p.IsDraft() || p.IsFuture() || p.IsExpired() {
// 			continue
// 		}
// 		title := p.Title
// 		path := p.Permalink()
// 		text, err := html2text.FromString(string(p.Content))
// 		if err != nil {
// 			log.Fatal("Could not convert Hugo content to text", err)
// 		}
// 		if *verbose {
// 			fmt.Println(path)
// 		 	fmt.Println("\t", title);
// 		 	//fmt.Println("\t", p.Summary);
// 			fmt.Println()
// 		}
// 		_, err = stmt.Exec(path, title, text, p.Summary)
// 	if err != nil {
// 		log.Fatal("Could not write page to DB", err)
// 	}
// 	}
// 	stmt.Close()
// }
