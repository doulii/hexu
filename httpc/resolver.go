package httpc

import (
	"context"
	"io"
	"net/http"
)

type Resolver interface {
	Resolve(ctx context.Context, method, url string, body io.Reader) (*http.Request, error)
}

type ResolveFunc func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error)

func (r ResolveFunc) Resolve(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return r(ctx, method, url, body)
}
