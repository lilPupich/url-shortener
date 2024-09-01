package delete

import (
	resp "URLShortener/internal/lib/api/response"
	"URLShortener/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteUrl(alias string) error
}

type Response struct {
	resp.Response
	Message string `json:"message"`
}

func New(log *slog.Logger, URLDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.url.delete.New"

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

		err := URLDeleter.DeleteUrl(alias)
		if err != nil {
			log.Error("failed to delete url", slog.String("alias", alias), sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("deleted alias", slog.String("alias", alias))
		successfulDeletion(w, r, "alias deleted successfully")

	}

}

func successfulDeletion(w http.ResponseWriter, r *http.Request, msg string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Message:  msg,
	})
}
