package tests

import (
	"URLShortener/internal/http-server/handlers/url/save"
	"URLShortener/internal/lib/api"
	"URLShortener/internal/lib/random"
	"net/http"
	"net/url"
	"path"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8082"
)

func TestUrlShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.New(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(5),
		}).
		WithBasicAuth("admin", "admin").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")
}

func TestUrlShortener_SaveRedirectRemove(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "valid url",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "invalid url",
			url:   "invalid",
			alias: random.NewRandomString(5),
			error: "field URL is not a valid url",
		},
		{
			name:  "empty alias",
			url:   gofakeit.URL(),
			alias: "",
		},
		{
			name:  "valid test100%",
			url:   "http://www.centralsynthesize.net/models",
			alias: "jPaVe",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.New(t, u.String())

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("admin", "admin").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().Equal(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			testRedirect(t, alias, tc.url)

			reqDel := e.DELETE("/"+path.Join("url", alias)).
				WithBasicAuth("admin", "admin").
				Expect().Status(200).
				JSON().Object()

			reqDel.Value("status").String().IsEqual("OK")

			testRedirectNotFound(t, alias)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	_, err := api.GetRedirect(u.String())
	require.ErrorIs(t, err, api.ErrInvalidStatusCode)
}
