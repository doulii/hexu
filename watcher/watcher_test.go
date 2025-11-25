package watcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestWatchString(t *testing.T) {
	type C struct {
		Addr Watcher[string]
	}
	data := `
addr: baidu.com
`
	cfg := &C{}
	err := yaml.Unmarshal([]byte(data), cfg)
	assert.Nil(t, err)
	assert.Equal(t, "baidu.com", cfg.Addr.Get())

	updated := false
	val := "baidu.com"
	cfg.Addr.Watch(func(v string) {
		// 回调会被调用两次，首次watch会调用一次，值变更会调用一次
		assert.Equal(t, val, v)
		updated = true
	})
	assert.True(t, updated)

	data2 := `
addr: google.com
`
	cfg2 := &C{}
	err = yaml.Unmarshal([]byte(data2), cfg2)
	assert.Nil(t, err)

	updated = false
	val = "google.com"
	cfg.Addr.Set(cfg2.Addr.Get()) // update
	assert.True(t, updated)
}

func TestWatchStruct(t *testing.T) {
	type V struct {
		Country string
		City    string
	}
	type C struct {
		Addr Watcher[*V]
	}
	data := `
addr:
  country: cn
  city: shanghai
`
	cfg := &C{}
	err := yaml.Unmarshal([]byte(data), cfg)
	assert.Nil(t, err)
	v := cfg.Addr.Get()
	expect := &V{Country: "cn", City: "shanghai"}
	assert.Equal(t, expect.Country, v.Country)
	assert.Equal(t, expect.City, v.City)

	updated := false
	cfg.Addr.Watch(func(v *V) {
		// 回调会被调用两次，首次watch会调用一次，值变更会调用一次
		assert.Equal(t, expect.Country, v.Country)
		assert.Equal(t, expect.City, v.City)
		updated = true
	})
	assert.True(t, updated)

	data2 := `
addr:
  country: us
  city: newyork
`
	cfg2 := &C{}
	err = yaml.Unmarshal([]byte(data2), cfg2)
	assert.Nil(t, err)

	updated = false
	expect = &V{Country: "us", City: "newyork"}
	cfg.Addr.Set(cfg2.Addr.Get()) // update
	assert.True(t, updated)
}
