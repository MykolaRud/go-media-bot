package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"media_bot/infrastructures"
	"media_bot/models"
	"media_bot/parsers"
	"media_bot/repositories"
	"time"
)

var (
	MySQLConfig = mysql.Config{
		User:   "root",
		Passwd: "diedie11",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "media_bot",
	}
	Repo *repositories.ArticleRepository
	Conn *sql.DB
)

func main() {

	initDB()
	initRepo()

	parsedArticleChannel := make(chan models.ParsedArticle)

	newsApiParser := parsers.NewsapiParser{}
	go func() {
		for {
			newsApiParser.ParseNew(parsedArticleChannel)
			time.Sleep(time.Minute)
		}
	}()

	redditParser := parsers.RedditParser{}

	go func() {
		for {
			redditParser.ParseNew(parsedArticleChannel)
			time.Sleep(time.Minute)
		}
	}()

	//go parser 1
	//go parser 2

	//read from channel
	//var parsedArticle models.ParsedArticle

	//for {
	//	select {
	//	case parsedArticle = <-parsedArticleChannel:
	//
	//		//fmt.Println("article grabbed ", parsedArticle.PublishedAt, ": ", parsedArticle.Title)
	//
	//		added := Repo.CheckAndAdd(parsedArticle)
	//		if added {
	//			fmt.Println("article received ", parsedArticle.PublishedAt, ": ", parsedArticle.Title)
	//		}
	//
	//	default:
	//		//fmt.Println("no message")
	//		time.Sleep(time.Millisecond * 300)
	//	}
	//}

	for parsedArticle := range parsedArticleChannel {
		added := Repo.CheckAndAdd(parsedArticle)
		if added {
			fmt.Println("article received ", parsedArticle.PublishedAt, ": ", parsedArticle.Title)
		}
	}

}

func initDB() {
	var err error

	Conn, err = sql.Open("mysql", MySQLConfig.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := Conn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func initRepo() {
	mysqlHandler := &infrastructures.MySQLHandler{}
	mysqlHandler.Conn = Conn

	Repo = &repositories.ArticleRepository{mysqlHandler}
}
