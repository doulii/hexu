package middleware

import (
	"net/http"
	"strings"

	"github.com/doulii/hexu/httpc/types"
	"github.com/doulii/hexu/log"
)

type Formatter interface {
	Format(r *http.Request) string
}

func NewLogger(f Formatter) func(next types.RequestFunc) types.RequestFunc {
	return func(next types.RequestFunc) types.RequestFunc {
		return func(req *http.Request) (*http.Response, error) {
			log.Warnf(req.Context(), "%s", f.Format(req))
			return next(req)
		}
	}
}

type CurlFormat struct {
	secretHeaders map[string]struct{}
}

func NewCurlFormat(secretHeaders []string) *CurlFormat {
	m := make(map[string]struct{}, len(secretHeaders))
	for _, h := range secretHeaders {
		m[http.CanonicalHeaderKey(h)] = struct{}{}
	}

	for _, h := range knownSecretHeaders {
		m[http.CanonicalHeaderKey(h)] = struct{}{}
	}
	return &CurlFormat{secretHeaders: m}
}

func (f *CurlFormat) Format(r *http.Request) string {
	sb := strings.Builder{}
	sb.WriteString("curl -X")
	sb.WriteString(r.Method)
	for k, vs := range r.Header {
		// k = http.CanonicalHeaderKey(k) // 写入时标准化
		if _, ok := f.secretHeaders[k]; ok {
			sb.WriteString(` -H"`)
			sb.WriteString(k)
			sb.WriteString(`: ***"`)
			continue
		}
		for _, v := range vs {
			sb.WriteString(` -H"`)
			sb.WriteString(k)
			sb.WriteString(`: `)
			sb.WriteString(v)
			sb.WriteString(`"`)
		}
	}
	sb.WriteString(` "`)
	sb.WriteString(r.URL.String())
	sb.WriteString(`"`)
	return sb.String()
}

var (
	knownSecretHeaders = []string{"Cookie"}
	Logger             = NewLogger(NewCurlFormat(nil))
)
