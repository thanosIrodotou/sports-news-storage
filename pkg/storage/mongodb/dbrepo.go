package mongodb

import (
	"context"

	"com.thanos/pkg/news"
)

//go:generate mockgen -source=dbrepo.go -destination=dbrepomock.go -package=mongodb

type DBRepo interface {
	GetArticleByID(context.Context, string) (Result, error)
	GetNews(context.Context) ([]Result, error)
	BulkInsert([]news.NewsArticle) (*BulkInsertResult, error)
}
