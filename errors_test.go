package gocaptcha

import "testing"

func TestErrorConstants(t *testing.T) {
    if ErrTokenNotFound == nil || ErrResponseNotFound == nil {
        t.Error("Error constants should be initialized")
    }
}
