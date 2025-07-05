package gocaptcha

import (
	"os"
	"testing"
	"time"
)

func TestErrorConstants(t *testing.T) {
    if ErrTokenNotFound == nil || ErrResponseNotFound == nil {
        t.Error("Error constants should be initialized")
    }
}

func TestParseURL(t *testing.T) {
    anchor := "api2/anchor?param1=foo&param2=bar"
    endpoint, params, err := parseURL(anchor)
    if err != nil {
        t.Fatalf("parseURL error: %v", err)
    }
    if endpoint != "api2" {
        t.Errorf("expected endpoint 'api2', got %s", endpoint)
    }
    if params != "param1=foo&param2=bar" {
        t.Errorf("unexpected params: %s", params)
    }
}

func TestRecaptchaV3(t *testing.T) {
	recaptcha, err := NewRecaptchaV3(os.Getenv("TEST_ANCHOR"), nil, 5*time.Second)
	if err != nil {
		t.Fatalf("NewRecaptchaV3 error: %v", err)
	}
	token, err := recaptcha.Solve()
	if err != nil {
		t.Fatalf("Solve error: %v", err)
	}
	if token == "" {
		t.Error("Solve returned an empty token")
	}
}