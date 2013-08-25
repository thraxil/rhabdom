package main

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/thraxil/paginate"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

func makeHandler(f func(http.ResponseWriter, *http.Request, *context),
	ctx *context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, ctx)
	}
}

type IndexResponse struct {
	Page     paginate.Page
	Posts    []Post
	SiteName string
}

func indexHandler(w http.ResponseWriter, r *http.Request, ctx *context) {
	index, _ := getIndex(ctx.PlainClient, ctx.PostCoder)
	var p = paginate.Paginator{ItemList: index, PerPage: 10}
	page := p.GetPage(r)
	iposts := page.Items()
	posts := make([]Post, len(iposts))
	for i, v := range iposts {
		posts[i] = v.(Post)
	}
	ir := IndexResponse{
		Page:     page,
		Posts:    posts,
		SiteName: SITE_NAME,
	}
	tmpl := getTemplate("index.html")
	tmpl.Execute(w, ir)
}

type AddResponse struct {
	YoutubeID string
	Title     string
	SiteName  string
}

func youtubeIDFromURL(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}
	q := u.Query()
	return q["v"][0]
}

func addHandler(w http.ResponseWriter, r *http.Request, ctx *context) {
	if r.Method == "POST" {
		title := r.PostFormValue("title")
		youtube_id := r.PostFormValue("youtube_id")
		NewPost(youtube_id, title, ctx)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		youtube_id := youtubeIDFromURL(r.FormValue("url"))
		title := string(r.FormValue("title"))
		r := "▶ "
		title = strings.Trim(title, r)
		tmpl := getTemplate("add.html")
		tmpl.Execute(w, AddResponse{YoutubeID: youtube_id,
			Title: title, SiteName: SITE_NAME})
	}
}

func atomHandler(w http.ResponseWriter, r *http.Request,
	ctx *context) {
	index, _ := getIndex(ctx.PlainClient, ctx.PostCoder)
	var p = paginate.Paginator{ItemList: index, PerPage: 20}
	page := p.GetPageNumber(1)
	iposts := page.Items()
	posts := make([]Post, len(iposts))
	for i, v := range iposts {
		posts[i] = v.(Post)
	}
	feed := &feeds.Feed{
		Title:       SITE_NAME,
		Link:        &feeds.Link{Href: FEED_LINK},
		Description: FEED_DESCRIPTION,
		Author:      &feeds.Author{FEED_AUTHOR_NAME, FEED_AUTHOR_EMAIL},
		Created:     posts[0].TimestampAsTime(),
	}

	feed.Items = []*feeds.Item{}
	for _, p := range posts {
		feed.Items = append(feed.Items,
			&feeds.Item{
				Title:       p.Title,
				Link:        &feeds.Link{Href: p.URL()},
				Description: p.Body(),
				Created:     p.TimestampAsTime(),
			})
	}
	atom, _ := feed.ToAtom()
	w.Header().Set("Content-Type", "application/atom+xml")
	fmt.Fprintf(w, atom)
}

type PostResponse struct {
	Post     Post
	SiteName string
}

func postHandler(w http.ResponseWriter, r *http.Request, ctx *context) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) < 3 {
		http.Error(w, "bad request", 400)
		return
	}
	id := parts[2]
	if id == "" {
		http.Error(w, "bad request", 400)
		return
	}
	var post Post
	_, err := ctx.PostCoder.FetchStruct(BUCKET, id, &post)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	ir := PostResponse{
		Post:     post,
		SiteName: SITE_NAME,
	}
	tmpl := getTemplate("post.html")
	tmpl.Execute(w, ir)
}

func getTemplate(filename string) *template.Template {
	var t = template.New("base.html")
	return template.Must(t.ParseFiles(
		filepath.Join(template_dir, "base.html"),
		filepath.Join(template_dir, filename),
	))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {}
