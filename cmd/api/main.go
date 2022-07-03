package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"com.thanos/pkg/api"
	"com.thanos/pkg/config"
	"com.thanos/pkg/logger"
	"com.thanos/pkg/storage/mongodb"
	"com.thanos/pkg/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var version string

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	if version != "" {
		cfg.APP.Version = version
		cfg.Logger.AppVersion = version
	}

	l := logger.NewLogger(cfg.Logger,
		logger.EnableReportCaller(),
		logger.SetLevel(cfg.Logger.LogLevel),
	)

	v, err := validator.New()
	if err != nil {
		l.WithError(err).Fatal("could not register validator translations")
	}

	mClient, err := mongodb.NewMongoClient(cfg.Mongo)
	if err != nil {
		l.WithError(err).Fatal("mongoDB connection unavailable")
	}

	// Select a database collection and inject it to repo
	db := mClient.Database(cfg.Mongo.Database)
	collection := db.Collection(cfg.Mongo.Collection)
	repo := mongodb.NewMongoRepo(collection, cfg.Mongo)

	if err = createIndexes(collection, cfg.Mongo.TTL); err != nil {
		l.WithError(err).Fatal("could not create db indexes")
	}

	a := api.NewAPI(
		api.NewJSONResponder(cfg.APP.Name, v.Translator),
		v,
		repo,
		cfg,
		l,
	)

	r := api.NewRouter(a, l)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		newNewsArticlesTicker := time.NewTicker(cfg.API.NewNewsArticlesFetchInterval).C

		for range newNewsArticlesTicker {
			articles := a.GetNewNewsArticles()
			if len(articles) > 0 {
				br, err := repo.BulkInsert(articles)
				if err != nil {
					l.WithError(err).Error("bulkInsert operation failed")
				}
				l.Infof("storing: %d articles, bulkInsert result: %+v", len(articles), br)
			}
		}
	}()

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// Start listening to incoming requests
	go func() {
		l.Infof("starting server at port: %d", cfg.Server.Port)
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			l.WithError(err).Fatal("server error")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// Gracefully Shutdown server
	l.Debug("app received termination signal, shutting down")
	if err := s.Shutdown(ctx); err != nil {
		l.WithError(err).Fatal("failed to gracefully shutdown http server")
	}

	cancel()
}

func createIndexes(coll *mongo.Collection, ttl time.Duration) error {
	indexOpts := options.Index()
	indexOpts.SetUnique(true)

	indexView := coll.Indexes()
	expireAfterSeconds := int32(ttl.Seconds())

	_, err := indexView.CreateMany(context.Background(),
		[]mongo.IndexModel{
			{
				Keys: bsonx.Doc{
					{Key: "articleID", Value: bsonx.Int32(1)},
				},
				Options: indexOpts,
			},
			{
				Keys: bsonx.Doc{
					{Key: "publishedAt", Value: bsonx.Int32(1)},
				},
				Options: &options.IndexOptions{ExpireAfterSeconds: &expireAfterSeconds},
			},
		})

	return err
}
