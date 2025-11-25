package watcher

import (
	"encoding/json"
	"sync"

	"gopkg.in/yaml.v3"
)

type Watcher[T any] struct {
	l        sync.RWMutex
	v        T
	watchers []func(T)
}

func (c *Watcher[T]) Get() T {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.v
}

func (c *Watcher[T]) Set(v T) {
	c.l.Lock()
	c.v = v
	wathers := c.watchers
	c.l.Unlock()
	for _, f := range wathers {
		f(v)
	}
}

// 订阅变更，会直接使用当前值调用回调函数
// 暂不支持取消订阅
func (c *Watcher[T]) Watch(f func(v T)) {
	c.l.Lock()
	c.watchers = append(c.watchers, f)
	v := c.v
	c.l.Unlock()
	f(v)
}

func (c *Watcher[T]) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(&c.v)
}

func (c *Watcher[T]) MarshalYAML() (interface{}, error) {
	return c.v, nil
}

// TODO: json test
func (c *Watcher[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.v)
}

func (c *Watcher[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.v)
}

// TODO: toml
