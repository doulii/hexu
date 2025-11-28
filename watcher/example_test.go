package watcher_test

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/doulii/hexu/watcher"
)

type Config struct {
	Weight watcher.Watcher[int] `json:"weight"`
}

type Service struct {
	cfg *Config
}

func NewService(cfg *Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Handle() string {
	if rand.Intn(100) < s.cfg.Weight.Get() {
		return "hit"
	}
	return "miss"
}

func ExampleWatcher() {
	cfg := &Config{}
	err := json.Unmarshal([]byte(`{"weight": 0}`), &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	svc := NewService(cfg)
	fmt.Println(svc.Handle())
	cfg.Weight.Set(100)
	fmt.Println(svc.Handle())
	// Output:
	// miss
	// hit
}
