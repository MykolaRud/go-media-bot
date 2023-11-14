package models

type NewsApiResponse struct {
	Status       string          `json:"status"`
	TotalResults int             `json:"totalResults"`
	Articles     []ParsedArticle `json:"articles"`
}
