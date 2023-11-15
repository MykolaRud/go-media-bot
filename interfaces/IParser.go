package interfaces

import "media_bot/models"

type IParser interface {
	ParseNew(c chan models.ParsedArticle)
	getFeed() []byte
}
