package main

import (
	"github.com/nu7hatch/gouuid"
	"log"
	"time"
)

type Post struct {
	YoutubeID string `riak:"youtube_id" json:"youtube_id"`
	Title     string `riak:"title" json:"title"`
	Timestamp string `riak:"timestamp" json:"timestamp"`
}

func NewPost(youtube_id, title string, ctx *context) {
	// Store Struct (uses coder)
	u4, err := uuid.NewV4()
	if err != nil {
		log.Println("error:", err)
		return
	}
	t := time.Now()
	data := Post{
		YoutubeID: youtube_id,
		Title:     title,
		Timestamp: t.Format(time.RFC3339),
	}
	if _, err := ctx.PostCoder.StoreStruct("test.rhabdom", u4.String(), &data); err != nil {
		log.Println(err.Error())
	}
	if err := ctx.PostCoder.LinkAdd("test.rhabdom", "index",
		"test.rhabdom", u4.String(), "post"); err != nil {
		log.Println(err.Error())
	}

	log.Println("saved struct")
}
