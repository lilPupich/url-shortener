package redirect

import (
	resp "URLShortener/internal/lib/api/response"
	"URLShortener/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, URLGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.url.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("empty alias")

			render.JSON(w, r, resp.Error("invalid alias"))

			return
		}

		resURL, err := URLGetter.GetUrl(alias)
		if err != nil {
			log.Error("failed to get url", slog.String("alias", alias), sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
