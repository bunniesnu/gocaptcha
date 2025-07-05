package gocaptcha

import (
	"fmt"
	"regexp"
)

// parseURL extracts endpoint and params from the anchor URL
func parseURL(anchor string) (endpoint, params string, err error) {
    re := regexp.MustCompile(`(?P<endpoint>(api2|enterprise))/anchor\?(?P<params>.*)`)
    m := re.FindStringSubmatch(anchor)
    if m == nil {
        return "", "", fmt.Errorf("invalid anchor URL")
    }
    // m[1] is endpoint, m[2] is params
    return m[1], m[3], nil
}