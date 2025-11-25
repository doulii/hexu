package httpc

import "github.com/doulii/hexu/httpc/types"

type Middleware func(next types.RequestFunc) types.RequestFunc
