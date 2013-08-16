package main

import (
	"github.com/mrb/riakpbc"
	"log"
)

func makeIndex(c *riakpbc.Client) {
	if _, err := c.StoreObject("test.rhabdom",
		"index", ""); err != nil {
		log.Println(err.Error())
	}
	log.Println("made index")
}

func getIndex(c *riakpbc.Client) (*riakpbc.RpbGetResp, error) {
	return c.FetchObject("test.rhabdom", "index")
}
