package gocaptcha

import "fmt"

// Type specifies proxy protocol
type Type string

const (
    HTTPS  Type = "https"
    SOCKS4 Type = "socks4"
    SOCKS5 Type = "socks5"
)

// Proxy holds proxy configuration
type Proxy struct {
    Type     Type
    Host     string
    Port     string
    Username string
    Password string
}

// URLs returns proxy URL strings for http and https
func (p Proxy) URLs() (httpURL, httpsURL string) {
    auth := ""
    if p.Username != "" || p.Password != "" {
        auth = fmt.Sprintf("%s:%s@", p.Username, p.Password)
    }
    scheme := string(p.Type)
    httpScheme := scheme
    if p.Type == HTTPS {
        httpScheme = "http"
    }
    return fmt.Sprintf("%s://%s%s:%s", httpScheme, auth, p.Host, p.Port),
        fmt.Sprintf("%s://%s%s:%s", scheme, auth, p.Host, p.Port)
}