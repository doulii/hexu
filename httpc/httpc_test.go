package httpc_test

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/doulii/hexu/httpc"
	"github.com/doulii/hexu/log"
)

func ExampleNew() {
	c := httpc.New(nil)
	resp, err := c.Get(context.Background(), "https://baidu.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	// Output:
	// 200 OK
}

func ExampleNewWithDefaultMiddleware() {
	log.Init(&log.Config{
		Sink: log.NewStdSink(os.Stdout),
	})
	c := httpc.NewWithDefaultMiddleware(nil)
	resp, err := c.Get(context.Background(), "https://baidu.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	// Output:
	// 200 OK
}

func ExampleOption() {
	c := httpc.New(nil)
	ctx := context.Background()
	resp, err := c.Get(ctx, "https://api.github.com/octocat", httpc.Header{
		"Accept":               "application/vnd.github+json",
		"X-GitHub-Api-Version": "2022-11-28",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}
