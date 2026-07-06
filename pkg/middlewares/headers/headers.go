package headers

import (
	"context"
	"net/http"
	"strings"
)

type headers struct {
	next http.Handler
	cfg  Config
}

type Config struct {
	CustomResponseHeaders map[string]string
}

func New(ctx context.Context, next http.Handler, cfg Config, name string) (http.Handler, error) {
	return &headers{
		next: next,
		cfg:  cfg,
	}, nil
}

func (h *headers) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.next.ServeHTTP(rw, req)
}

func (h *headers) ModifyResponse(res *http.Response) error {
	// Save original Set-Cookie headers to prevent them from being collapsed
	cookies := res.Header.Values("Set-Cookie")
	hasCookies := len(cookies) > 0

	for k, v := range h.cfg.CustomResponseHeaders {
		if strings.EqualFold(k, "Set-Cookie") {
			continue
		}
		if v == "" {
			res.Header.Del(k)
		} else {
			res.Header.Set(k, v)
		}
	}

	if hasCookies {
		res.Header.Del("Set-Cookie")
		for _, cookie := range cookies {
			res.Header.Add("Set-Cookie", cookie)
		}
	}

	return nil
}
