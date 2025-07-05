package gocaptcha

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Session wraps http.Client
type Session struct {
    client *http.Client
}

// NewSession creates a new Session with optional proxy and timeout
func NewSession(proxy *Proxy, timeout time.Duration) (*Session, error) {
    transport := &http.Transport{}
    if proxy != nil {
        _, httpsURL := proxy.URLs()
        proxyURL, err := url.Parse(httpsURL)
        if err != nil {
            return nil, err
        }
        transport.Proxy = http.ProxyURL(proxyURL)
    }
    client := &http.Client{
        Transport: transport,
        Timeout:   timeout,
    }
    return &Session{client: client}, nil
}

// sendRequest sends GET or POST to the given endpoint with params/data
func (s *Session) sendRequest(method, endpoint, params, data string) (*http.Response, error) {
    fullURL := fmt.Sprintf(BaseURL, endpoint)
    var req *http.Request
    var err error
    paramValues, err := url.ParseQuery(params)
    if err != nil {
        return nil, fmt.Errorf("invalid params: %w", err)
    }
    if method == http.MethodPost {
        req, err = http.NewRequest(http.MethodPost, fullURL, strings.NewReader(data))
        if err != nil {
            return nil, err
        }
        req.Header.Set("Content-Type", ContentType)
    } else {
        req, err = http.NewRequest(http.MethodGet, fullURL, nil)
        if err != nil {
            return nil, err
        }
    }
    q := req.URL.Query()
    for key, values := range paramValues {
        for _, value := range values {
            q.Add(key, value)
        }
    }
    req.URL.RawQuery = q.Encode()
    return s.client.Do(req)
}