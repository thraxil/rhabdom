package main

import (
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
	Page  paginate.Page
	Posts []Post
}

func indexHandler(w http.ResponseWriter, r *http.Request, ctx *context) {
	index, _ := getIndex(ctx.PlainClient, ctx.PostCoder)
	var p = paginate.Paginator{ItemList: index, PerPage: 20}
	page := p.GetPage(r)
	iposts := page.Items()
	posts := make([]Post, len(iposts))
	for i, v := range iposts {
		posts[i] = v.(Post)
	}
	ir := IndexResponse{
		Page:  page,
		Posts: posts,
	}
	tmpl := getTemplate("index.html")
	tmpl.Execute(w, ir)
}

type AddResponse struct {
	YoutubeID string
	Title string
}

func youtubeIDFromURL(rawurl string) string {
	u, err  := url.Parse(rawurl)
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
		tmpl.Execute(w, AddResponse{YoutubeID: youtube_id, Title: title})
	}
}

func getTemplate(filename string) *template.Template {
	var t = template.New("base.html")
	return template.Must(t.ParseFiles(
		filepath.Join(template_dir, "base.html"),
		filepath.Join(template_dir, filename),
	))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {}
