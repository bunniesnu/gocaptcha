package gocaptcha

import "errors"

var (
    ErrTokenNotFound    = errors.New("recaptcha token not found")
    ErrResponseNotFound = errors.New("recaptcha response not found")
)