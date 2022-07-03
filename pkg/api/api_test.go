package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.thanos/pkg/api"
	"com.thanos/pkg/config"
	"com.thanos/pkg/logger"
	"com.thanos/pkg/news"
	"com.thanos/pkg/storage/mongodb"
	"com.thanos/pkg/validator"
	"github.com/golang/mock/gomock"
)

func TestAPI_GetAllArticles(t *testing.T) {
	testCases := []struct {
		description          string
		newsArticles         []mongodb.Result
		expectedStatus       int
		expectedNewsArticles int
		expectedError        error
	}{
		{
			description: "should respond with 200 and a list of news articles",
			newsArticles: []mongodb.Result{
				{
					ID:        "62c3063d04b7e4c1864ff552",
					ArticleID: "645168",
					Data: news.Data{
						Id:          "645168",
						TeamId:      "",
						OptaMatchId: nil,
						Title:       "Brentford FC advertise for Technology Support Technician",
						Type:        []string{},
						Teaser:      t,
						Content:     "",
						Url:         "https://www.brentfordfc.com/news/2022/july/brentford-fc-advertise-for-technology-support-technician/",
						ImageUrl:    "",
						GalleryUrls: nil,
						VideoUrl:    nil,
						Published:   "2022-07-04 13:00:00",
					},
				},
				{
					ID:        "62c3063d04b7e4c1864ff553",
					ArticleID: "645150",
					Data: news.Data{
						Id:          "645150",
						TeamId:      "",
						OptaMatchId: nil,
						Title:       "Club supports fans heading to tournament",
						Type:        []string{},
						Teaser:      t,
						Content:     "",
						Url:         "https://www.brentfordfc.com/news/2022/july/brentford-fc-represented-at-worldnet-2022/",
						ImageUrl:    "",
						GalleryUrls: nil,
						VideoUrl:    nil,
						Published:   "2022-07-04 11:00:00",
					},
				},
			},
			expectedStatus:       http.StatusOK,
			expectedNewsArticles: 2,
		},
		{
			description:          "should respond with 500 and an error response",
			newsArticles:         []mongodb.Result{},
			expectedStatus:       http.StatusInternalServerError,
			expectedNewsArticles: 0,
			expectedError:        errors.New("storage error"),
		},
	}

	cfg, err := config.New("../../")
	if err != nil {
		t.Fatal(err)
	}

	log := logger.NewLogger(cfg.Logger, logger.DisableOutput())

	v, err := validator.New()
	if err != nil {
		t.Fatal(err)
	}

	responder := api.NewJSONResponder(cfg.APP.Name, v.Translator)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.description, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/articles", nil)

			// Mock API dependencies
			ctrl := gomock.NewController(t)

			dbrepo := mongodb.NewMockDBRepo(ctrl)
			if tc.expectedError != nil {
				dbrepo.EXPECT().GetNews().Return(tc.newsArticles, tc.expectedError)
			} else {
				dbrepo.EXPECT().GetNews().Return(tc.newsArticles, nil)
			}

			a := api.NewAPI(responder, v, dbrepo, cfg, log)

			recorder := httptest.NewRecorder()
			a.ErrorWrapper(a.GetAllArticles).ServeHTTP(recorder, request)

			var resp api.Response
			if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
				t.Fatal(err)
			}

			if recorder.Code != tc.expectedStatus {
				t.Fatalf("expected to get status %d, got %d", tc.expectedStatus, recorder.Code)
			}
		})
	}
}
