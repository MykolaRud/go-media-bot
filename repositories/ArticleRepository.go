package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"media_bot/interfaces"
	"media_bot/models"
	"time"
)

type ArticleRepository struct {
	interfaces.IDbHandler
}

func (repo *ArticleRepository) CheckAndAdd(article models.ParsedArticle) bool {
	if repo.ArticleExists(article) {
		return false
	} else {
		_, err := repo.AddArticle(article)
		if err != nil {
			fmt.Errorf("error adding article %s", err)
		}

		return true
	}
}

func (repo *ArticleRepository) ArticleExists(article models.ParsedArticle) bool {
	var existingArticle models.Article

	row := repo.QueryRow("SELECT id FROM articles WHERE url = ?", article.Url[:min(250, len(article.Url))])
	err := row.Scan(&existingArticle)
	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func (repo *ArticleRepository) AddArticle(article models.ParsedArticle) (int64, error) {

	url := article.Url[:min(250, len(article.Url))]
	title := article.Title[:min(250, len(article.Title))]

	Result, err := repo.Execute("INSERT INTO articles (source_id, url, title, created_at) VALUES (?, ?, ?, ?)",
		article.SourceId,
		url,
		title,
		time.Now(),
	)

	if err != nil {
		log.Fatal(err)

		return 0, err
	}

	lastInsertId, err := Result.LastInsertId()
	if err != nil {
		log.Fatal(err)

		return 0, err
	}

	return lastInsertId, nil
}
