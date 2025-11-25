package httpc

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/doulii/hexu/log"
)

// api级别使用[]byte作为body
// 如果是超大的body，应当直接使用DoRequest

type Input interface {
	MarshalBody() ([]byte, error)
	ContentType() string
}

type Output interface {
	UnmarshalBody([]byte) error
}

func (c *client) DoAPI(ctx context.Context, method, url string, input Input, output Output, options ...Option) (err error) {
	data, err := input.MarshalBody()
	if err != nil {
		return fmt.Errorf("marshal input: %w", err)
	}
	resp, err := c.Do(ctx, method, url, bytes.NewReader(data), append(options, Header{
		"Content-Type": input.ContentType(),
	})...)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer CloseBody(ctx, resp.Body)
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	err = output.UnmarshalBody(data)
	if err != nil {
		return fmt.Errorf("unmarshal response body: %w, data=%s", err, data)
	}
	return
}

func CloseBody(ctx context.Context, body io.Closer) {
	err := body.Close()
	if err != nil {
		log.Errorf(ctx, "close body error %v", err)
	}
}
