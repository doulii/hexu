package watcher

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate_Simple(t *testing.T) {
	o := Watcher[int]{}
	n := Watcher[int]{}
	o.Set(1)
	n.Set(2)
	assert.Equal(t, 1, o.Get())
	assert.Equal(t, 2, n.Get())
	Update(&o, &n)
	assert.Equal(t, 2, o.Get())
	assert.Equal(t, 2, n.Get())
}

func TestUpdate_Field(t *testing.T) {
	type Config struct {
		Weight Watcher[int] `json:"weight"`
	}

	o := &Config{}
	err := json.Unmarshal([]byte(`{"weight": 10}`), o)
	assert.Nil(t, err)
	assert.Equal(t, 10, o.Weight.Get())

	n := &Config{}
	err = json.Unmarshal([]byte(`{"weight": 20}`), n)
	assert.Nil(t, err)
	assert.Equal(t, 20, n.Weight.Get())

	Update(o, n)
	assert.Equal(t, 20, o.Weight.Get())
	assert.Equal(t, 20, n.Weight.Get())
}

// unrecommended
func TestUpdate_PointerField(t *testing.T) {
	type Config struct {
		Weight *Watcher[int] `json:"weight"`
	}

	o := &Config{}
	err := json.Unmarshal([]byte(`{"weight": 10}`), o)
	assert.Nil(t, err)
	assert.Equal(t, 10, o.Weight.Get())

	n := &Config{}
	err = json.Unmarshal([]byte(`{"weight": 20}`), n)
	assert.Nil(t, err)
	assert.Equal(t, 20, n.Weight.Get())

	Update(o, n)
	assert.Equal(t, 20, o.Weight.Get())
	assert.Equal(t, 20, n.Weight.Get())
}

// unrecommended
// nil *Watcher can not be updated
func TestUpdate_NilPointerField(t *testing.T) {
	type Config struct {
		Weight *Watcher[int] `json:"weight"`
	}

	c1 := &Config{}
	err := json.Unmarshal([]byte(`{}`), c1)
	assert.Nil(t, err)
	assert.Nil(t, c1.Weight)

	c2 := &Config{}
	err = json.Unmarshal([]byte(`{"weight": 20}`), c2)
	assert.Nil(t, err)
	assert.Equal(t, 20, c2.Weight.Get())

	Update(c1, c2)
	assert.Nil(t, c1.Weight) // update failed
	assert.Equal(t, 20, c2.Weight.Get())

	Update(c2, c1)
	assert.Nil(t, c1.Weight)
	assert.Equal(t, 20, c2.Weight.Get()) // not updated
}

func TestUpdate_Nested(t *testing.T) {
	type Inner struct {
		Weight Watcher[int]
	}
	type Outer struct {
		Inner1 Inner
		Inner2 *Inner
	}
	o := &Outer{
		Inner1: Inner{
			Weight: Watcher[int]{},
		},
		Inner2: &Inner{
			Weight: Watcher[int]{},
		},
	}
	o.Inner1.Weight.Set(1)
	o.Inner2.Weight.Set(2)
	n := &Outer{
		Inner1: Inner{
			Weight: Watcher[int]{},
		},
		Inner2: &Inner{
			Weight: Watcher[int]{},
		},
	}
	n.Inner1.Weight.Set(3)
	n.Inner2.Weight.Set(4)

	assert.Equal(t, 1, o.Inner1.Weight.Get())
	assert.Equal(t, 2, o.Inner2.Weight.Get())
	assert.Equal(t, 3, n.Inner1.Weight.Get())
	assert.Equal(t, 4, n.Inner2.Weight.Get())
	Update(o, n)
	assert.Equal(t, 3, o.Inner1.Weight.Get())
	assert.Equal(t, 4, o.Inner2.Weight.Get())
	assert.Equal(t, 3, n.Inner1.Weight.Get())
	assert.Equal(t, 4, n.Inner2.Weight.Get())
}

func TestUpdate_Map(t *testing.T) {
	type Config map[string]*Watcher[int]

	c1 := make(Config)
	err := json.Unmarshal([]byte(`{"a":10,"b":20}`), &c1)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(c1))
	assert.Equal(t, 10, c1["a"].Get())
	assert.Equal(t, 20, c1["b"].Get())

	c2 := make(Config)
	err = json.Unmarshal([]byte(`{"a":30,"c":40}`), &c2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(c2))
	assert.Equal(t, 30, c2["a"].Get())
	assert.Equal(t, 40, c2["c"].Get())

	Update(c1, c2)
	assert.Equal(t, 2, len(c1))
	assert.Equal(t, 30, c1["a"].Get())
	assert.Equal(t, 20, c1["b"].Get())
	assert.Equal(t, 2, len(c2))
	assert.Equal(t, 30, c2["a"].Get())
	assert.Equal(t, 40, c2["c"].Get())
}
