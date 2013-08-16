package main

type Post struct {
	YoutubeID string `riak:"youtube_id" json:"youtube_id"`
	Title     string `riak:"title" json:"title"`
	Timestamp string `riak:"timestamp" json:"timestamp"`
}
