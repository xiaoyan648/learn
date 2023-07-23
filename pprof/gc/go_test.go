package gc

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Cache = *wrapper

type wrapper struct {
	*cache
}

type cache struct {
	content   string
	stop      chan struct{}
	onStopped func()
}

func newCache() *cache {
	return &cache{
		content: "some thing",
		stop:    make(chan struct{}),
	}
}

func NewCache() Cache {
	w := &wrapper{
		cache: newCache(),
	}
	go w.cache.run()
	runtime.SetFinalizer(w, stopShardedJanitor)
	return w
}

func stopShardedJanitor(w *wrapper) {
	w.stop <- struct{}{}
}

func (w *wrapper) Stop() {
	w.cache.Stop()
}

func (c *cache) run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// do some thing
		case <-c.stop:
			if c.onStopped != nil {
				c.onStopped()
			}
			return
		}
	}
}

func (c *cache) Stop() {
	close(c.stop)
}

func TestFinalizer(t *testing.T) {
	s := assert.New(t)

	w := NewCache()
	var cnt int = 0
	stopped := make(chan struct{})
	w.onStopped = func() {
		cnt++
		close(stopped)
	}

	s.Equal(0, cnt)

	w = nil

	runtime.GC()

	select {
	case <-stopped:
	case <-time.After(10 * time.Second):
		t.Fail()
	}

	s.Equal(1, cnt)
}
