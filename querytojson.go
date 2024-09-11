// Package querytojson is a traefik plugin for converting URL queries into JSON field in request body.
package querytojson

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// compile-time interface assertion
var _ http.Handler = (*QueryToJSON)(nil)

// Config the plugin configuration.
type Config struct {
	TargetJSONField string `json:"targetQuery,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		TargetJSONField: "queryStringParameters", // AWS API Gateway...
	}
}

type QueryToJSON struct {
	next http.Handler
	name string
	cfg  *Config
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &QueryToJSON{
		next: next,
		name: name,
		cfg:  config,
	}, nil
}

func (a *QueryToJSON) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	jsonBody := make(map[string]any)
	if req.Body != nil {
		if err := json.NewDecoder(req.Body).Decode(&jsonBody); err != nil {
			// make an empty JSON if request body is invalid JSON
			jsonBody = map[string]any{}
		}
	}

	jsonFields := make(map[string]string)
	queries := req.URL.Query()
	for k, query := range queries {
		// support only first field for now. empty, or mulitple queries will be ignored.
		// e.g. http://example.com/index?product=1&color=red&order&hobbies=programming&hobbies=sports
		// product/color is valid query. order will be ignored completely.
		// hobbies will be used, but only last one will be used.
		// TODO: Support multiple query fields
		for _, q := range query {
			jsonFields[k] = q
		}
	}
	jsonBody[a.cfg.TargetJSONField] = jsonFields

	jb, err := json.Marshal(jsonBody)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Body = io.NopCloser(strings.NewReader(string(jb)))
	req.ContentLength = int64(len(jb))

	a.next.ServeHTTP(rw, req)
}
