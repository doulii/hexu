package types

import "net/http"

type RequestFunc func(req *http.Request) (*http.Response, error)
