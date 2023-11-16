package parsers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"media_bot/interfaces"
	"media_bot/models"
	"net/http"
	"slices"
)

type NewsapiParser struct {
	interfaces.IParser
}

func (parser NewsapiParser) ParseNew(c chan models.ParsedArticle) {
	feed := parser.getFeed()
	newsApiResponse := parser.parseFeed(feed)
	parser.pushFeedData(newsApiResponse, c)
}

func (parser NewsapiParser) getFeed() []byte {

	ApiKey := "f88bdcacb34e4b2e8b3f5bd84724f741"

	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=Football&from=2023-11-07&sortBy=publishedAt&apiKey=%s", ApiKey)
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func (parser NewsapiParser) parseFeed(feed []byte) models.NewsApiResponse {
	var newsApiResponse models.NewsApiResponse
	json.Unmarshal(feed, &newsApiResponse)

	return newsApiResponse
}

func (parser NewsapiParser) pushFeedData(newsApiResponse models.NewsApiResponse, c chan models.ParsedArticle) {
	ignoreArticleTitles := []string{"[Removed]"}

	source := models.Source{}

	for i := range newsApiResponse.Articles {
		if slices.Index(ignoreArticleTitles, newsApiResponse.Articles[i].Title) == -1 {
			article := newsApiResponse.Articles[i]
			article.SourceId = source.SourceIdNewsApi()

			c <- article
		}
	}
}
