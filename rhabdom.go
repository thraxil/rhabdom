package main

import (
	"flag"
	"github.com/mrb/riakpbc"
	"github.com/stvp/go-toml-config"
	"log"
	"net/http"
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
