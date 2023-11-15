package parsers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"media_bot/interfaces"
	"media_bot/models"
	"net/http"
	"strings"
	"time"
)

type RedditParser struct {
	interfaces.IParser
	UrlPrefix string
}

func (parser RedditParser) ParseNew(c chan models.ParsedArticle) {
	parser.UrlPrefix = "https://www.reddit.com"

	feed := parser.getFeed()
	parsedArticles := parser.parseFeed(feed)
	parser.pushFeedData(parsedArticles, c)
}

func (parser RedditParser) getFeed() []byte {
	//url := "https://www.reddit.com/r/popular/new/"
	url := parser.UrlPrefix + "/t/nfl/"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func (parser RedditParser) parseFeed(feed []byte) []models.ParsedArticle {
	var parsedArticles []models.ParsedArticle
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(feed[:])))
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("shreddit-feed shreddit-post").Each(func(i int, s *goquery.Selection) {
		a := s.Find("a[slot=\"full-post-link\"]")
		title, titleExists := a.Attr("aria-label")
		url, _ := a.Attr("href")
		url = parser.UrlPrefix + url
		createdAt, _ := s.Attr("created-timestamp")

		if !titleExists {
			return
		}
		layout := "2006-01-02T15:04:05.000000+0000"
		createdAtTime, err := time.Parse(layout, createdAt)
		if err != nil {
			fmt.Println("error parsing time ", err)
		}

		parsedArticle := models.ParsedArticle{
			Title:       title,
			Url:         url,
			PublishedAt: createdAtTime,
		}

		parsedArticles = append(parsedArticles, parsedArticle)

	})

	return parsedArticles
}

func (parser RedditParser) pushFeedData(parsedArticles []models.ParsedArticle, c chan models.ParsedArticle) {
	source := models.Source{}

	for _, el := range parsedArticles {

		article := el
		article.SourceId = source.SourceIdReddit()

		//fmt.Printf("%t %v \n", article, article)
		//fmt.Println(article)

		c <- article
	}

	return
}
