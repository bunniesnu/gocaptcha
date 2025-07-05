package gocaptcha

import (
	"testing"
)

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