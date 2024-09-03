package redirect_test

import (
	"URLShortener/internal/http-server/handlers/redirect"
	"URLShortener/internal/http-server/handlers/redirect/mocks"
	"URLShortener/internal/lib/api"
	"URLShortener/internal/lib/logger/handlers/slogdiscard"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {

	testTable := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "success",
			alias: "testtt",
			url:   "http://test.com",
		},
	}

	for _, tc := range testTable {
		t.Run(
			tc.name,
			func(t *testing.T) {
				urlGetterMock := mocks.NewURLGetter(t)

				if tc.respError == "" || tc.mockError != nil {
					urlGetterMock.On("GetUrl", tc.alias).
						Return(tc.url, tc.mockError).Once()
				}

				r := chi.NewRouter()
				r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

				ts := httptest.NewServer(r)
				defer ts.Close()

				redirectedToUrl, err := api.GetRedirect(ts.URL + "/" + tc.alias)
				require.NoError(t, err)

				require.Equal(t, tc.url, redirectedToUrl)
			},
		)
	}
}
