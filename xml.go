package main

import (
	"encoding/xml"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// Sitemap
type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url    `xml:"url"`
}

type Url struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// Atom feed
type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Xmlns   string   `xml:"xmlns,attr"`
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Link    []Link   `xml:"link"`
	Entries []Entry  `xml:"entry"`
}

type Link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length int64  `xml:"length,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
}

type Content struct {
	Type string `xml:"type,attr,omitempty"`
	Body string `xml:",chardata"`
}

type Entry struct {
	Title     string     `xml:"title"`
	Link      Link       `xml:"link"`
	ID        string     `xml:"id"`
	Updated   string     `xml:"updated"`
	Summary   string     `xml:"summary,omitempty"`
	Content   Content    `xml:"content,omitempty"`
	Enclosure *Enclosure `xml:"enclosure,omitempty"`
}

func create_sitemap() error {
	var urls []Url

	fpath := filepath.Join(projectPath, "dist")

	filepath.Walk(fpath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".html" {
			rel, _ := filepath.Rel(fpath, path)
			url := filepath.Join(config.SiteUrl, rel)
			urls = append(urls, Url{
				Loc:     url,
				LastMod: time.Now().Format("2006-01-02"),
			})
		}
		return nil
	})

	sitemap := UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}

	data, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(projectPath, "dist", "sitemap.xml"), data, 0644); err != nil {
		return err
	}

	return nil
}

func create_feed() error {
	fpath, err := url.JoinPath(config.SiteUrl, "atom.xml")
	if err != nil {
		return err
	}

	feed := Feed{
		Xmlns:   "http://www.w3.org/2005/Atom",
		Title:   config.SiteTitle,
		ID:      config.SiteUrl,
		Updated: time.Now().Format(time.RFC3339),
		Link: []Link{
			{Href: fpath, Rel: "self"},
			{Href: config.SiteUrl},
		},
		Entries: postsFeed,
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(projectPath, "dist", "atom.xml"), data, 0644); err != nil {
		return err
	}

	return nil
}
