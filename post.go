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
	Id        string `riak:"id" json:"id"`
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
		Id:        u4.String(),
	}
	if _, err := ctx.PostCoder.StoreStruct("test.rhabdom", u4.String(), &data); err != nil {
		log.Println(err.Error())
	}
	if err := ctx.PostCoder.LinkAdd("test.rhabdom", "index",
		"test.rhabdom", u4.String(), "post"); err != nil {
		log.Println(err.Error())
	}
}

func (p Post) Body() string {
	return `<iframe class="youtube-player" type="text/html"
  width="640" height="385" 
    src="http://www.youtube.com/embed/` + p.YoutubeID + `"
    allowfullscreen frameborder="0">`
}

func (p Post) URL() string {
	return URL_BASE + "/post/" + p.Id + "/"
}

func (p Post) TimestampAsTime() time.Time {
	created, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", p.Timestamp)
	if err != nil {
		return time.Now()
	}
	return created
}
