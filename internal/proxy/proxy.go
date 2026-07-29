package proxy

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewProxy(upstreamURL string, resolver func(string) string, logger *slog.Logger) *httputil.ReverseProxy {
	target, err := url.Parse(upstreamURL)
	if err != nil {
		logger.Error("Invalid upstream URL", "url", upstreamURL)
		return nil
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)
			pr.Out.Header.Del("Accept-Encoding")
		},

		ModifyResponse: func(res *http.Response) error {
			reqPath := res.Request.URL.Path

			if res.StatusCode == http.StatusOK &&
				strings.Contains(reqPath, "/store/pos-version/") &&
				res.Request.Method == http.MethodGet {

				bodyBytes, err := io.ReadAll(res.Body)
				if err != nil {
					return err
				}
				res.Body = io.NopCloser(bytes.NewReader(bodyBytes))

				logger.Info("Body", "body", string(bodyBytes))

				posVersion := strings.TrimSpace(string(bodyBytes))
				posVersion = strings.Trim(posVersion, `"`)

				bundleID := resolver(posVersion)
				if bundleID != "" {
					logger.Info("PosVersion intercepted, injecting cookie", "posVersion", posVersion, "bundle", bundleID)
					cookie := &http.Cookie{
						Name:     "pos_version",
						Value:    bundleID,
						Path:     "/",
						HttpOnly: false,
						SameSite: http.SameSiteStrictMode,
					}
					res.Header.Add("Set-Cookie", cookie.String())
				} else {
					logger.Warn("PosVersion intercepted but NO MATCH in mapping", "posVersion", posVersion)
				}
			}
			return nil
		},
	}

	return proxy
}
