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
		port       = config.String("port", "9999")
		media_dir  = config.String("media_dir", "media")
		t_dir      = config.String("template_dir", "templates")
		riak_nodes = config.String("riak_nodes", "")
		site_name  = config.String("site_name", "Rhabdom")
	)
	config.Parse(configFile)
	template_dir = *t_dir
	SITE_NAME = *site_name
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
