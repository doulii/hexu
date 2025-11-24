package httpc

import (
	"context"
	"io"
	"net/http"
)

type Client interface {
	Do(req *http.Request) (resp *http.Response, err error)

	Head(ctx context.Context, url string, options ...Option) (resp *http.Response, err error)
	Get(ctx context.Context, url string, options ...Option) (resp *http.Response, err error)
	Delete(ctx context.Context, url string, options ...Option) (resp *http.Response, err error)

	Post(ctx context.Context, url string, body io.Reader, options ...Option) (resp *http.Response, err error)
	Put(ctx context.Context, url string, body io.Reader, options ...Option) (resp *http.Response, err error)

	// DoAPI(ctx context.Context, methdo, url string, input Input, output Output, options ...Option) (err error)
	// PostJson(ctx context.Context, url string, input, output interface{}, options ...Option) (err error)
}

type Config struct {
	UserAgent string
	Client    *http.Client
	Resolver  Resolver
}

type client struct {
	userAgent   string
	client      *http.Client
	middlewares []Middleware

	do RequestFunc

	resolver Resolver
}

/*
1. 底层net/http Do
2. middleware
3. Request Builder
4. json/xml API
*/
func New(cfg *Config, middlewares ...Middleware) (c *client) {
	c = &client{
		userAgent:   defaultUserAgent,
		client:      http.DefaultClient,
		middlewares: middlewares,
		resolver:    ResolveFunc(http.NewRequestWithContext),
	}
	c.mergeConfig(cfg)
	do := c.client.Do
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		do = c.middlewares[i](do)
	}
	c.do = do
	return
}

func (c *client) mergeConfig(cfg *Config) {
	if cfg == nil {
		return
	}
	if cfg.UserAgent != "" {
		c.userAgent = cfg.UserAgent
	}
	if cfg.Client != nil {
		c.client = cfg.Client
	}
	if cfg.Resolver != nil {
		c.resolver = cfg.Resolver
	}
}

func (c *client) DoRequest(req *http.Request) (resp *http.Response, err error) {
	if req.Header.Get(UserAgentHeader) == "" {
		req.Header.Set(UserAgentHeader, c.userAgent)
	}
	return c.do(req)
}

func (c *client) Do(ctx context.Context, method, url string, body io.Reader, options ...Option) (resp *http.Response, err error) {
	req, err := c.resolver.Resolve(ctx, method, url, body)
	if err != nil {
		return
	}
	for _, opt := range options {
		if err = opt.apply(req); err != nil {
			return
		}
	}
	return c.DoRequest(req)
}

func (c *client) Head(ctx context.Context, url string, options ...Option) (resp *http.Response, err error) {
	return c.Do(ctx, http.MethodHead, url, nil, options...)
}

func (c *client) Get(ctx context.Context, url string, options ...Option) (resp *http.Response, err error) {
	return c.Do(ctx, http.MethodGet, url, nil, options...)
}

func (c *client) Post(ctx context.Context, url string, body io.Reader, options ...Option) (resp *http.Response, err error) {
	return c.Do(ctx, http.MethodPost, url, body, options...)
}

func (c *client) Put(ctx context.Context, url string, body io.Reader, options ...Option) (resp *http.Response, err error) {
	return c.Do(ctx, http.MethodPut, url, body, options...)
}

func (c *client) Patch(ctx context.Context, url string, body io.Reader, options ...Option) (resp *http.Response, err error) {
	return c.Do(ctx, http.MethodPatch, url, body, options...)
}

func (c *client) Delete(ctx context.Context, url string, options ...Option) (resp *http.Response, err error) {
	return c.Do(ctx, http.MethodDelete, url, nil, options...)
}
