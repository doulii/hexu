package httpc

import "net/http"

type RequestFunc func(req *http.Request) (*http.Response, error)
type Middleware func(next RequestFunc) RequestFunc
