package main

import (
	"github.com/mrb/riakpbc"
	"log"
)

type MainIndex struct {
	Obj       *riakpbc.RpbGetResp
	PostCoder *riakpbc.Client
}

func makeIndex(c *riakpbc.Client) {
	if _, err := c.StoreObject(BUCKET,
		"index", ""); err != nil {
		log.Println(err.Error())
	}
	log.Println("made index")
}

func getIndex(c *riakpbc.Client, pc *riakpbc.Client) (*MainIndex, error) {
	o, err := c.FetchObject(BUCKET, "index")
	if err != nil {
		return nil, err
	}
	return &MainIndex{Obj: o, PostCoder: pc}, nil
}

func (m MainIndex) TotalItems() int {
	return len(m.Obj.GetContent()[0].GetLinks())
}

func (m MainIndex) ItemRange(offset, count int) []interface{} {
	total := m.TotalItems()
	posts := make([]interface{}, count)

	for i := 0; i < count; i++ {
		key := m.Obj.GetContent()[0].GetLinks()[total-(offset+i+1)].Key
		lpost := &Post{}
		if _, err := m.PostCoder.FetchStruct(
			BUCKET, string(key), lpost); err != nil {
			log.Println(err.Error())
		}

		posts[i] = *lpost
	}
	return posts
}
