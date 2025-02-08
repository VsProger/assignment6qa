package handlers

import (
	"bytes"
	"forum/app"
	mock "forum/internal/repo/mocks"
	"forum/internal/service"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const (
	sessionNameInCookie = "session"
	sessionCookieValue  = "anythingHereWouldWork"
)

type TestServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T) *TestServer {
	var buff bytes.Buffer

	logger := log.New(&buff, "", 0)

	templateCache, err := app.NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	app := app.New(logger, logger, templateCache)
	repo := mock.NewMockRepo(t)
	serv := service.New(repo)

	hand := New(serv, app)

	ts := httptest.NewServer(hand.Routes())

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestServer{ts}
}

func (ts *TestServer) get(t *testing.T, url string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + url)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *TestServer) postForm(t *testing.T, url string, form url.Values) (int, http.Header, string) {
	req, err := http.NewRequest("POST", ts.URL+url, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.AddCookie(&http.Cookie{
		Name:  sessionNameInCookie,
		Value: sessionCookieValue,
	})

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return res.StatusCode, res.Header, string(body)
}
