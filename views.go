package main

import (
	"github.com/mrb/riakpbc"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func makeHandler(f func(http.ResponseWriter, *http.Request, *context),
	ctx *context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, ctx)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request, ctx *context) {
	index, _ := getIndex(ctx.PlainClient)
	for i, k := range index.GetContent()[0].GetLinks() {
		log.Println(i, string(k.Key))
	}

	tmpl := getTemplate("index.html")
	tmpl.Execute(w, nil)
}

func addHandler(w http.ResponseWriter, r *http.Request, ctx *context) {
	addPost(ctx.PostCoder, ctx.PlainClient)
}

func addPost(postCoder, plainClient *riakpbc.Client) {
	// Store Struct (uses coder)
	data := Post{
		YoutubeID: "p43CfAVgX2U",
		Title:     "Siobhan Donaghy - Ghosts",
		Timestamp: "2012-08-12 22:03:00",
	}
	if _, err := postCoder.StoreStruct("test.rhabdom", "testpost", &data); err != nil {
		log.Println(err.Error())
	}
	if err := postCoder.LinkAdd("test.rhabdom", "index",
		"test.rhabdom", "testpost", "post"); err != nil {
		log.Println(err.Error())
	}

	log.Println("saved struct")
	// Fetch Struct (uses coder)
	out := &Post{}
	if _, err := postCoder.FetchStruct("test.rhabdom", "testpost", out); err != nil {
		log.Println(err.Error())
	}
	log.Println(out.Title)
}

func getTemplate(filename string) *template.Template {
	var t = template.New("base.html")
	return template.Must(t.ParseFiles(
		filepath.Join(template_dir, "base.html"),
		filepath.Join(template_dir, filename),
	))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {}
