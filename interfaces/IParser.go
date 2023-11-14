package interfaces

import "media_bot/models"

type IParser interface {
	ParseNew(c chan models.ParsedArticle)
	getFeed() []byte
	parseFeed(feed []byte) models.NewsApiResponse
	pushFeedData(newsApiResponse models.NewsApiResponse, c chan models.ParsedArticle)
}
