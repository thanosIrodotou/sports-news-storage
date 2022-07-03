package api

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"com.thanos/pkg/config"
	"com.thanos/pkg/logger"
	"com.thanos/pkg/news"
	"com.thanos/pkg/storage/mongodb"
	"com.thanos/pkg/validator"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

const ISO8601 = "2006-01-02T15:04:05.000Z"

// Handler custom handler signature to be able to return errors & handle them centrally
type Handler func(w http.ResponseWriter, r *http.Request) error

type API struct {
	Responder
	validate   *validator.Validator
	repository mongodb.DBRepo
	cfg        *config.Config
	log        *logger.Logger
}

// NewAPI creates a new API
func NewAPI(
	r Responder,
	v *validator.Validator,
	repo mongodb.DBRepo,
	c *config.Config,
	l *logger.Logger,
) *API {
	return &API{
		Responder:  r,
		validate:   v,
		repository: repo,
		cfg:        c,
		log:        l,
	}
}

// GetNewNewsArticles this is called periodically by a goroutine and its
// purpose is to fetch news articles periodically
// Ideally this functionality would be on a separate module and not on the api
func (a *API) GetNewNewsArticles() []news.NewsArticle {
	uri := fmt.Sprintf("%s%d", a.cfg.API.GetLatestNewsArticlesUrl, a.cfg.API.NewsArticlesPerCall)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		a.log.WithError(err).Error("could not create latest news articles request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.log.WithError(err).Error("could not retrieve latest news articles")
	}
	defer resp.Body.Close()

	var newListInformation news.NewListInformation
	bytez, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		a.log.WithError(err).Error("could not read response bytes")
	}

	err = xml.NewDecoder(bytes.NewReader(bytez)).Decode(&newListInformation)
	if err != nil {
		a.log.WithError(err).Error("could not decode response bytes")
	}

	newsItems := newListInformation.NewsletterNewsItems.NewsletterNewsItem
	newsArticles := make([]news.NewsArticle, 0, len(newsItems))

	for i := range newsItems {
		ni := newsItems[i]
		newsArticle := news.NewsArticle{
			Data: news.Data{
				Id:          ni.NewsArticleID,
				Published:   ni.PublishDate,
				Title:       ni.Title,
				OptaMatchId: ni.OptaMatchId,
				Url:         ni.ArticleURL,
			},
		}
		newsArticles = append(newsArticles, newsArticle)
	}

	return newsArticles
}

// GetAllArticles retrieve all articles
func (a *API) GetAllArticles(w http.ResponseWriter, r *http.Request) error {
	newsArticles, err := a.repository.GetNews(r.Context())
	if err != nil {
		return a.RespondError(r.Context(), w, err)
	}

	return a.Respond(
		r.Context(),
		w,
		Response{
			Status: "success",
			Data:   newsArticles,
			Metadata: news.Metadata{
				CreatedAt:  time.Now().UTC().Format(ISO8601),
				Sort:       "-published",
				TotalItems: len(newsArticles),
			},
		},
		http.StatusOK,
	)
}

// GetArticleByID retrieve an article by its unique ID
// TODO: unimplement
func (a *API) GetArticleByID(w http.ResponseWriter, r *http.Request) error {
	var id string
	if sid := chi.URLParam(r, "id"); sid != "" {
		_, err := strconv.ParseInt(sid, 10, 32)
		if err != nil {
			return ErrBadRequest
		}
		id = sid
	}

	newsArticle, err := a.repository.GetArticleByID(r.Context(), id)
	if err != nil {
		return a.RespondError(r.Context(), w, err)
	}

	return a.Respond(
		r.Context(),
		w,
		Response{
			Status: "success",
			Data:   newsArticle,
			Metadata: news.Metadata{
				CreatedAt:  time.Now().UTC().Format(ISO8601),
				Sort:       "-published",
				TotalItems: 1,
			},
		},
		http.StatusOK,
	)
}

// Health check
// TODO: health response should check whether the db connection is stable.
// If not, api's health should update to reflect that
func (a *API) Health(w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte{})
	return err
}

// Version returns build version
func (a *API) Version(w http.ResponseWriter, r *http.Request) error {
	resp := VersionResponse{Version: a.cfg.APP.Version}
	return a.Respond(r.Context(), w, resp, http.StatusOK)
}

// ErrorWrapper wrap custom handler signatures to handle & log their errors
func (a *API) ErrorWrapper(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logFields := logrus.Fields{
			"url.path": r.URL.RequestURI(),
		}

		statusCode := http.StatusInternalServerError
		errorType := "errInternalServerError"
		errorMessage := ""
		message := ""

		if err := h(w, r); err != nil {
			if err, ok := err.(Error); ok {
				message = err.message
				errorMessage = err.message
				statusCode = err.statusCode
				errorType = err.Code
			}

			// Prepare log entry and log
			logFields["http.response.status_code"] = http.StatusOK
			logFields["error.code"] = statusCode
			logFields["error.type"] = errorType
			if errorMessage != "" {
				logFields["error.message"] = errorMessage
			}

			_, file, line, ok := runtime.Caller(1)
			if ok {
				logFields["log.origin.file.name"] = file
				logFields["log.origin.file.line"] = line
			}

			a.log.WithFields(logFields).Error(message)

			// Respond to caller
			if err := a.RespondError(r.Context(), w, err); err != nil {
				a.log.WithFields(logFields).WithError(err).Error("failed to responsd")
			}
		} else {
			logFields["http.response.status_code"] = http.StatusOK
			a.log.WithFields(logFields).Debug("transaction ended")
		}
	}
}
