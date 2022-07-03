package api

import (
	"com.thanos/pkg/logger"
	"github.com/go-chi/chi"
)

type Router struct {
	*chi.Mux
}

// NewRouter returns a new handler with all registered routes
func NewRouter(api *API, log *logger.Logger) *Router {
	rt := Router{Mux: chi.NewRouter()}

	rt.Route("/v1", func(r chi.Router) {
		r.Get("/articles", api.ErrorWrapper(api.GetAllArticles))
		r.Get("/article/{id}", api.ErrorWrapper(api.GetArticleByID))
	})

	rt.Get("/version", api.ErrorWrapper(api.Version))
	rt.Get("/health", api.ErrorWrapper(api.Health))

	return &rt
}
