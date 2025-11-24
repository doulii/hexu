package httpc

import (
	"net/http"
	"net/url"
)

type Option interface {
	apply(*http.Request) error
}

type Header map[string]string

func (h Header) apply(r *http.Request) error {
	for k, v := range h {
		r.Header.Set(k, v)
	}
	return nil
}

type HttpHeader http.Header

func (h HttpHeader) apply(r *http.Request) error {
	for k, vs := range h {
		r.Header.Del(k)
		for _, v := range vs {
			r.Header.Add(k, v)
		}
	}
	return nil
}

type Query map[string]string

func (q Query) apply(r *http.Request) error {
	query := r.URL.Query()
	for k, v := range q {
		query.Set(k, v)
	}
	r.URL.RawQuery = query.Encode()
	return nil
}

type HttpQuery url.Values

func (q HttpQuery) apply(r *http.Request) error {
	query := r.URL.Query()
	for k, v := range q {
		for _, e := range v {
			query.Add(k, e)
		}
	}
	r.URL.RawQuery = query.Encode()
	return nil
}

type CustomOption func(*http.Request) error

func (c CustomOption) apply(r *http.Request) error {
	return c(r)
}
