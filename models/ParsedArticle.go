package models

import "time"

type ParsedArticle struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	PublishedAt time.Time `json:"PublishedAt"`
	SourceId    int64
}
