package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	template_mocks "github.com/alienix2/sensor_visualization/cmd/web/mocks"
	mocks "github.com/alienix2/sensor_visualization/internal/models/mocks"
)

func newTestApplication() *application {
	templateCache := template_mocks.FakeTemplateCache()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		infoLog:        log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog:       log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		accounts:       &mocks.MockAccountModel{},
		topics:         &mocks.MockTopicModel{},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) post(t *testing.T, urlPath string, data io.Reader) (int, http.Header, string) {
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/x-www-form-urlencoded", data)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}

func LoadAndSaveMock(session *scs.SessionManager, key string, value int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return session.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session.Put(r.Context(), key, value)
			next.ServeHTTP(w, r)
		}))
	}
}
