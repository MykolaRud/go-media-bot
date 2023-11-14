package models

import "time"

type Article struct {
	SourceId  int
	URL       string
	Title     string
	CreatedAt time.Time
}
