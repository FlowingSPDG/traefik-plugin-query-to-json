package querytojson_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	querytojson "github.com/FlowingSPDG/traefik-plugin-query-to-json"
)

func TestQueryToJSON(t *testing.T) {
	cfg := querytojson.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

	handler, err := querytojson.New(ctx, next, cfg, "query-to-json")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost?hoge=fuga", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	_ = body
	// TODO: assert w/o external libraries
}
