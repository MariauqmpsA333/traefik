package service

import (
	"net/http"
	"net/http/httputil"
)

type ResponseModifier interface {
	ModifyResponse(*http.Response) error
}

func NewProxy(transport http.RoundTripper, responseModifiers []ResponseModifier) *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{
		Transport: transport,
	}

	proxy.ModifyResponse = func(res *http.Response) error {
		// Save original Set-Cookie headers to prevent them from being collapsed by modifiers
		cookies := res.Header.Values("Set-Cookie")
		hasCookies := len(cookies) > 0

		for _, modifier := range responseModifiers {
			err := modifier.ModifyResponse(res)
			if err != nil {
				return err
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

	return proxy
}
