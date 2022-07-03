package mongodb_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"com.thanos/pkg/config"
	"com.thanos/pkg/news"
	"com.thanos/pkg/storage/mongodb"
	"github.com/sirupsen/logrus"
)

var repository *mongodb.Repository

func TestMain(m *testing.M) {
	var err error

	cfg, err := config.New()
	if err != nil {
		fmt.Printf("could not load configuration, failed to start tests. %v", err)
		os.Exit(1)
	}

	cfg.Mongo.Host = "localhost"
	cfg.Mongo.Port = 27100
	cfg.Mongo.TTL = 10 * time.Minute

	log := logrus.New()
	log.Out = ioutil.Discard

	mClient, err := mongodb.NewMongoClient(cfg.Mongo)
	if err != nil {
		fmt.Printf("could not initialize mongodb client, %v", err)
		os.Exit(1)
	}

	db := mClient.Database(cfg.Mongo.Database)
	collection := db.Collection(cfg.Mongo.Collection)
	repository = mongodb.NewMongoRepo(collection, cfg.Mongo)

	// set seed so that random is semi-predictable
	rand.Seed(54)

	os.Exit(m.Run())
}

func TestMongoDBRepo_GetNews(t *testing.T) {
	randomArticle := newArticle("1234")

	_, err := repository.BulkInsert(
		[]news.NewsArticle{randomArticle},
	)
	if err != nil {
		t.Fatal(err)
	}

	n, err := repository.GetNews(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(n) == 0 {
		t.Fatal("news articles returned should be > 0")
	}

	// TODO: perform deep equality check here
	if randomArticle.Data.Content != n[0].Data.Content {
		t.Fatal("news article contents should match")
	}

	// TODO: all randomly generated newsArticles need to be cleaned at the end of the test to avoid
	// errors on consecutive runs.
}

func TestMongoDBRepo_GetArticleByID(t *testing.T) {
	id := "4321"
	randomArticle := newArticle(id)

	_, err := repository.BulkInsert(
		[]news.NewsArticle{randomArticle},
	)
	if err != nil {
		t.Fatal(err)
	}

	n, err := repository.GetArticleByID(context.TODO(), id)
	if err != nil {
		t.Fatal(err)
	}

	if n.ArticleID != id {
		t.Fatalf("could not find article with id: %s, got: %s", id, n.ArticleID)
	}

	// TODO: perform deep equality check here
	if randomArticle.Data.Content != n.Data.Content {
		t.Fatal("news article contents should match")
	}
}

func newArticle(id string) news.NewsArticle {
	return news.NewsArticle{
		Data: news.Data{
			Id:          id,
			TeamId:      randString(24),
			OptaMatchId: rand.Int(),
			Title:       randString(24),
			Type:        []string{},
			Teaser:      nil,
			Content:     randString(24),
			Url:         randString(24),
			ImageUrl:    randString(24),
			GalleryUrls: nil,
			VideoUrl:    nil,
			Published:   time.Now().Format(mongodb.DATE_TIME_FORMAT),
		},
		Metadata: news.Metadata{},
		Status:   "success",
	}
}

func randString(n int) string {
	alphabet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(b)
}
