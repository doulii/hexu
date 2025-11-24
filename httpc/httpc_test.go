package httpc

import (
	"context"
	"fmt"
)

func ExampleNew() {
	c := New(nil)
	resp, err := c.Get(context.Background(), "https://baidu.com")
	fmt.Println(err)
	fmt.Println(resp.Status)
	// Output:
	// <nil>
	// 200 OK
}
