package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/dig"
	"log"
	"media_bot/infrastructures"
	"media_bot/interfaces"
	"media_bot/models"
	"media_bot/parsers"
	"media_bot/repositories"
	"time"
)

var (
	MySQLConfig = mysql.Config{
		User:   "rud",
		Passwd: "diedie11",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "media_bot",
	}
	Repo *repositories.ArticleRepository
)

func main() {
	container := dig.New()
	container.Provide(initDBConnection)
	container.Provide(initDBHandler)
	container.Invoke(initRepo)

	parsedArticleChannel := make(chan models.ParsedArticle)

	parsers := []interfaces.IParser{
		parsers.NewsapiParser{},
		parsers.RedditParser{},
	}

	for _, parser := range parsers {
		go func() {
			for {
				parser.ParseNew(parsedArticleChannel)
				time.Sleep(time.Minute)
			}
		}()
	}

	for parsedArticle := range parsedArticleChannel {
		added := Repo.CheckAndAdd(parsedArticle)
		if added {
			fmt.Println("article received ", parsedArticle.PublishedAt, ": ", parsedArticle.Title)
		}
	}

}

func initDBConnection() *sql.DB {
	var err error

	Conn, err := sql.Open("mysql", MySQLConfig.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := Conn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	return Conn
}

func initDBHandler() interfaces.IDbHandler {
	return &infrastructures.MySQLHandler{}
}

func initRepo(dbHandler interfaces.IDbHandler, Conn *sql.DB) {
	dbHandler.SetConn(Conn)

	Repo = &repositories.ArticleRepository{dbHandler}
}
