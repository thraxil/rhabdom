package main

import (
	"flag"
	"github.com/mrb/riakpbc"
	"github.com/stvp/go-toml-config"
	"log"
	"net/http"
	"strings"
)

var SITE_NAME = "Rhabdom"
var URL_BASE = "http://localhost:9999"
var FEED_LINK = "http://localhost:9999/atom.xml"
var FEED_DESCRIPTION = "video blog"
var FEED_AUTHOR_NAME = "anders pearson"
var FEED_AUTHOR_EMAIL = "anders@columbia.edu"
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
		port              = config.String("port", "9999")
		media_dir         = config.String("media_dir", "media")
		t_dir             = config.String("template_dir", "templates")
		riak_nodes        = config.String("riak_nodes", "")
		site_name         = config.String("site_name", "Rhabdom")
		url_base          = config.String("url_base", "http://localhost:9999")
		feed_link         = config.String("feed_link", "http://localhost:9999/atom.xml")
		feed_description  = config.String("feed_description", "video blog")
		feed_author_name  = config.String("feed_author_name", "anders")
		feed_author_email = config.String("feed_author_email", "anders@columbia.edu")
	)
	config.Parse(configFile)
	template_dir = *t_dir
	SITE_NAME = *site_name
	URL_BASE = *url_base
	FEED_LINK = *feed_link
	FEED_DESCRIPTION = *feed_description
	FEED_AUTHOR_NAME = *feed_author_name
	FEED_AUTHOR_EMAIL = *feed_author_email
	postCoder, plainClient, err := connect(strings.Split(*riak_nodes, ","))
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
	http.HandleFunc("/atom.xml", makeHandler(atomHandler, &ctx))
	http.HandleFunc("/add/", makeHandler(addHandler, &ctx))
	http.Handle("/media/", http.StripPrefix("/media/",
		http.FileServer(http.Dir(*media_dir))))
	http.ListenAndServe(":"+*port, nil)

	postCoder.Close()
	plainClient.Close()
}

func connect(nodes []string) (*riakpbc.Client, *riakpbc.Client, error) {
	coder := riakpbc.NewCoder("json", riakpbc.JsonMarshaller,
		riakpbc.JsonUnmarshaller)
	postCoder := riakpbc.NewClientWithCoder(nodes, coder)
	if err := postCoder.Dial(); err != nil {
		log.Print(err.Error())
		return nil, nil, err
	}
	if _, err := postCoder.SetClientId("rhabdom"); err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}

	plainClient := riakpbc.NewClient(nodes)
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
