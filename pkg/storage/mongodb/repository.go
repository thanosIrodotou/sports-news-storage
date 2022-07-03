package mongodb

import (
	"context"
	"fmt"
	"time"

	"com.thanos/pkg/config"
	"com.thanos/pkg/news"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DATE_TIME_FORMAT = "2006-01-02 15:04:05"

// Repository mongo struct
type Repository struct {
	articlesCollection *mongo.Collection
	cfg                config.Mongo
}

// BulkInsertResult represents the result of a bulk insert operation
type BulkInsertResult struct {
	InsertedCount int64
	UpsertedCount int64
}

// NewMongoRepo creates a new Mongo repository
func NewMongoRepo(n *mongo.Collection, c config.Mongo) *Repository {
	return &Repository{
		articlesCollection: n,
		cfg:                c,
	}
}

// Result represents a mongodb result
type Result struct {
	ID        string    `json:"-" bson:"_id"`
	ArticleID string    `json:"-" bson:"articleID"`
	Data      news.Data `json:"data"`
}

func (r Repository) GetArticleByID(ctx context.Context, id string) (newsArticle Result, err error) {
	result := r.articlesCollection.FindOne(ctx, bson.D{{Key: "articleID", Value: id}})

	if err = result.Decode(&newsArticle); err != nil {
		if err == mongo.ErrNoDocuments {
			return newsArticle, fmt.Errorf("article id (%s) does not exist: %w", id, err)
		}
	}

	return newsArticle, err
}

// GetNews returns a list of all newArticles
func (r Repository) GetNews(ctx context.Context) (newsArticles []Result, err error) {
	sortStage := bson.D{{"$sort", bson.D{{"published", -1}}}}
	cursor, err := r.articlesCollection.Aggregate(
		context.Background(),
		mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "publishedAt", Value: bson.D{
						{Key: "$gte", Value: time.Now().Add(-r.cfg.TTL).UTC()},
					}},
				}},
			},
			sortStage,
		},
		options.Aggregate(),
	)

	if err != nil {
		return newsArticles, err
	}

	for cursor.Next(ctx) {
		article := Result{}

		err = cursor.Decode(&article)
		if err != nil {
			return newsArticles, err
		}

		newsArticles = append(newsArticles, article)
	}

	return newsArticles, err
}

// BulkInsert inserts an array of NewsArticles using upsert to avoid dupes
func (r Repository) BulkInsert(news []news.NewsArticle) (*BulkInsertResult, error) {
	// Update records in any order
	bulkWriteOpts := options.BulkWrite()
	bulkWriteOpts.SetOrdered(false)

	models := make([]mongo.WriteModel, len(news))

	for i, n := range news {
		bulkModel := mongo.NewUpdateOneModel()

		dt, err := time.Parse(DATE_TIME_FORMAT, n.Data.Published)
		if err != nil {
			return nil, err
		}

		n.Metadata.CreatedAt = time.Now().Format(time.RFC3339)

		model := bulkModel.SetFilter(bson.D{
			{Key: "articleID", Value: n.Data.Id},
		}).SetUpdate(bson.D{
			{Key: "$set", Value: n},
			{Key: "$set", Value: bson.D{
				{Key: "publishedAt", Value: dt},
			}},
		}).SetUpsert(true)

		models[i] = model
	}

	res, err := r.articlesCollection.BulkWrite(context.Background(), models, bulkWriteOpts)
	result := newBulkWriteResult(res)

	return &result, err
}

func newBulkWriteResult(bwr *mongo.BulkWriteResult) BulkInsertResult {
	if bwr == nil {
		return BulkInsertResult{
			InsertedCount: 0,
			UpsertedCount: 0,
		}
	}

	return BulkInsertResult{
		InsertedCount: bwr.InsertedCount,
		UpsertedCount: bwr.UpsertedCount,
	}
}
