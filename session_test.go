// file: session_test.go
package gocaptcha

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// rewriteTransport is a RoundTripper that sends every request
// to tsURL, preserving the original path and query.
type rewriteTransport struct {
	tsURL string
}

func (rt rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// build a new URL: tsURL + original path + "?" + raw query
	target := rt.tsURL + req.URL.Path
	if req.URL.RawQuery != "" {
		target += "?" + req.URL.RawQuery
	}
	// construct a new request against the test server
	newReq, err := http.NewRequest(req.Method, target, req.Body)
	if err != nil {
		return nil, err
	}
	// copy over headers
	for k, vv := range req.Header {
		for _, v := range vv {
			newReq.Header.Add(k, v)
		}
	}
	return http.DefaultTransport.RoundTrip(newReq)
}

func TestSession_sendRequest(t *testing.T) {
	// 1) spin up test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Expect the query param "params"
			if got := r.URL.Query().Get("params"); got != "myParam" {
				t.Errorf("GET: expected params=myParam, got %q", got)
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "GET_OK")

		case http.MethodPost:
			// Expect the query param "k"
			if got := r.URL.Query().Get("k"); got != "postParam" {
				t.Errorf("POST: expected k=postParam, got %q", got)
			}
			// Expect the body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("POST: read body error: %v", err)
			}
			if string(body) != "foo=bar" {
				t.Errorf("POST: expected body 'foo=bar', got %q", body)
			}
			// Expect Content-Type
			if ct := r.Header.Get("Content-Type"); ct != ContentType {
				t.Errorf("POST: expected Content-Type %q, got %q", ContentType, ct)
			}
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, "POST_OK")

		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	}))
	defer ts.Close()

	// 2) build a Session whose transport rewrites to ts.URL
	client := &http.Client{
		Timeout:   2 * time.Second,
		Transport: rewriteTransport{tsURL: ts.URL},
	}
	sess := &Session{client: client}

	// 3) run subtests

	t.Run("GET", func(t *testing.T) {
		resp, err := sess.sendRequest(http.MethodGet, "ignored-endpoint", "myParam", "")
		if err != nil {
			t.Fatalf("sendRequest(GET) error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		payload, _ := io.ReadAll(resp.Body)
		if string(payload) != "GET_OK" {
			t.Errorf("expected body GET_OK, got %q", payload)
		}
	})

	t.Run("POST", func(t *testing.T) {
		resp, err := sess.sendRequest(http.MethodPost, "ignored-endpoint", "postParam", "foo=bar")
		if err != nil {
			t.Fatalf("sendRequest(POST) error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}
		payload, _ := io.ReadAll(resp.Body)
		if string(payload) != "POST_OK" {
			t.Errorf("expected body POST_OK, got %q", payload)
		}
	})
}
