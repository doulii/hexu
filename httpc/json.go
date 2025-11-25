package httpc

import (
	"context"
	"encoding/json"
	"net/http"
)

type JsonInput struct{ Val interface{} }

func (i *JsonInput) MarshalBody() ([]byte, error) {
	return json.Marshal(i.Val)
}

func (i *JsonInput) ContentType() string {
	return "application/json"
}

type JsonOutput struct{ Val interface{} }

func (o *JsonOutput) UnmarshalBody(data []byte) error {
	return json.Unmarshal(data, o.Val)
}

type JsonBody struct{ Val interface{} }

func (j *JsonBody) ContentType() string {
	return "application/json"
}

func (j *JsonBody) MarshalBody() ([]byte, error) {
	return json.Marshal(j.Val)
}

func (j *JsonBody) UnmarshalBody(data []byte) error {
	return json.Unmarshal(data, j.Val)
}

func (c *client) DoJsonAPI(ctx context.Context, method, url string, input, output interface{}, options ...Option) (err error) {
	return c.DoAPI(ctx, method, url, &JsonBody{input}, &JsonBody{output}, options...)
	// return c.DoAPI(ctx, method, url, &JsonInput{input}, &JsonOutput{output}, options...)
}

func (c *client) PostJson(ctx context.Context, url string, input, output interface{}, options ...Option) (err error) {
	return c.DoJsonAPI(ctx, http.MethodPost, url, input, output, options...)
}

func (c *client) PutJson(ctx context.Context, url string, input, output interface{}, options ...Option) (err error) {
	return c.DoJsonAPI(ctx, http.MethodPut, url, input, output, options...)
}
