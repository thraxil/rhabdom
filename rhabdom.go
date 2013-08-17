package main

import (
	"flag"
	"fmt"
	"github.com/mrb/riakpbc"
	"github.com/stvp/go-toml-config"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var RIAK_NODES = []string{"127.0.0.1:10017",
	"127.0.0.0:10027", "127.0.0.0:10037"}
var template_dir = "templates"

type context struct {
	PostCoder   *riakpbc.Client
	PlainClient *riakpbc.Client
}

func main() {
	var configFile string
	var makeindices string
	flag.StringVar(&configFile, "config", "./dev.conf", "TOML config file")
	flag.StringVar(&makeindices, "makeindices", "", "make index")
	flag.Parse()
	var (
		port      = config.String("port", "9999")
		media_dir = config.String("media_dir", "media")
		t_dir     = config.String("template_dir", "templates")
	)
	config.Parse(configFile)
	template_dir = *t_dir

	postCoder, plainClient, err := connect(RIAK_NODES)
	if err != nil {
		log.Println("couldn't connect to riak")
		return
	}
	if makeindices != "" {
		makeIndex(plainClient)
	}
	var ctx = context{
		PostCoder:   postCoder,
		PlainClient: plainClient,
	}
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", makeHandler(indexHandler, &ctx))
	http.HandleFunc("/add/", makeHandler(addHandler, &ctx))
	http.Handle("/media/", http.StripPrefix("/media/",
		http.FileServer(http.Dir(*media_dir))))
	http.ListenAndServe(":"+*port, nil)

	postCoder.Close()
	plainClient.Close()
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {}

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

func connect(nodes []string) (*riakpbc.Client, *riakpbc.Client, error) {
	coder := riakpbc.NewCoder("json", riakpbc.JsonMarshaller,
		riakpbc.JsonUnmarshaller)
	postCoder := riakpbc.NewClientWithCoder(RIAK_NODES, coder)
	if err := postCoder.Dial(); err != nil {
		log.Print(err.Error())
		return nil, nil, err
	}
	if _, err := postCoder.SetClientId("rhabdom"); err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}

	plainClient := riakpbc.NewClient(RIAK_NODES)
	if err := plainClient.Dial(); err != nil {
		log.Print(err.Error())
		return nil, nil, err
	}
	if _, err := plainClient.SetClientId("rhabdom-plain"); err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}
	return postCoder, plainClient, nil
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

	fmt.Println("saved struct")
	// Fetch Struct (uses coder)
	out := &Post{}
	if _, err := postCoder.FetchStruct("test.rhabdom", "testpost", out); err != nil {
		log.Println(err.Error())
	}
	fmt.Println(out.Title)
}

func getTemplate(filename string) *template.Template {
	var t = template.New("base.html")
	return template.Must(t.ParseFiles(
		filepath.Join(template_dir, "base.html"),
		filepath.Join(template_dir, filename),
	))
}
