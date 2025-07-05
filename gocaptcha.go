// recapcha_v3.go
package gocaptcha

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// RecaptchaV3 encapsulates the anchor URL, proxy/timeout config and a session.
type RecaptchaV3 struct {
    anchorURL string
    session   *Session
}

// NewRecaptchaV3 creates the recaptcha object and underlying HTTP session.
func NewRecaptchaV3(anchorURL string, proxy *Proxy, timeout time.Duration) (*RecaptchaV3, error) {
    sess, err := NewSession(proxy, timeout)
    if err != nil {
        return nil, err
    }
    return &RecaptchaV3{
        anchorURL: anchorURL,
        session:   sess,
    }, nil
}

// Solve runs the full two‑step exchange and returns the final token.
func (r *RecaptchaV3) Solve() (string, error) {
    // 1. parse anchor URL
    endpoint, rawParams, err := parseURL(r.anchorURL)
    if err != nil {
        return "", err
    }

    // 2. GET the recaptcha-token
    token, err := r.getRecaptchaToken(endpoint, rawParams)
    if err != nil {
        return "", err
    }

    // 3. split rawParams (“v=…&k=…&co=…”)
    parts := strings.Split(rawParams, "&")
    paramMap := make(map[string]string, len(parts))
    for _, p := range parts {
        kv := strings.SplitN(p, "=", 2)
        if len(kv) == 2 {
            paramMap[kv[0]] = kv[1]
        }
    }

    // 4. build POST body
    postBody := fmt.Sprintf(POST_DATA,
        paramMap["v"],    // version
        token,            // token from GET
        paramMap["k"],    // site key
        paramMap["co"],   // origin
    )

    // 5. POST to reload and extract the final response
    respToken, err := r.getRecaptchaResponse(endpoint, paramMap["k"], postBody)
    if err != nil {
        return "", err
    }

    return respToken, nil
}

func (r *RecaptchaV3) getRecaptchaToken(endpoint, params string) (string, error) {
    urlPath := fmt.Sprintf("%s/anchor", endpoint)
    resp, err := r.session.sendRequest(http.MethodGet, urlPath, params, "")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    re := regexp.MustCompile(`"recaptcha-token" value="(.*?)"`)
    matches := re.FindStringSubmatch(string(body))
    if len(matches) < 2 {
        return "", ErrTokenNotFound
    }
    return matches[1], nil
}

func (r *RecaptchaV3) getRecaptchaResponse(endpoint, siteKey, data string) (string, error) {
    urlPath := fmt.Sprintf("%s/reload", endpoint)
    // we pass siteKey as the param "k"
    resp, err := r.session.sendRequest(http.MethodPost, urlPath, siteKey, data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    re := regexp.MustCompile(`"rresp","(.*?)"`)
    matches := re.FindStringSubmatch(string(body))
    if len(matches) < 2 {
        return "", ErrResponseNotFound
    }
    return matches[1], nil
}
