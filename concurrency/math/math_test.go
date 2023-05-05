package math

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//=============== 并发安全测试 ==============//

// 自定义rand.NewSource() 存在并发问题
// WARNING: DATA RACE
func TestMathCustomRead(t *testing.T) {
	var wg sync.WaitGroup
	// rand.NewSource() 存在并发问题
	var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var tid [16]byte
			_, _ = randSource.Read(tid[:])
		}()
	}
	wg.Wait()
}

// 默认的全局rand是加锁的，无并发问题
// go test -timeout 30s -run ^TestMathDefaultRead$ github.com/xiaoyan648/learn/concurrency/math -count=1 -race
func TestMathDefaultRead(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var tid [16]byte
			_, _ = rand.Read(tid[:])
		}()
	}
	wg.Wait()
}

//=============== 性能测试 ==============//

// 方案1.使用全局rand（加全局锁）
// 无并发 BenchmarkMathRead-8   	51075153	        23.11 ns/op	       0 B/op	       0 allocs/op
// 并发 BenchmarkMathReadDefault-8   	 4055941	       296.0 ns/op	      16 B/op	       1 allocs/op
func BenchmarkMathReadDefault(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var tid [16]byte
			_, _ = rand.Read(tid[:])
		}()
	}
	wg.Wait()
}

// 方案2.每次new一个新rand对象
// 并发 BenchmarkMathReadNew-8   	  531456	      2232 ns/op	    5417 B/op	       2 allocs/op
func BenchmarkMathReadNew(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var tid [16]byte
			n := rand.New(rand.NewSource(time.Now().UnixNano()))
			_, _ = n.Read(tid[:])
		}()
	}
	wg.Wait()
}

// 方案3. 方案1的优化，使用分片锁减少锁冲突
// 并发 BenchmarkMathShardingMutex-8   	 4769182	       249.4 ns/op	      24 B/op	       1 allocs/op
func BenchmarkMathShardingMutex(b *testing.B) {
	sr := NewShardedRander(300)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var tid [16]byte
			_, _ = sr.Read(tid[:])
		}()
	}
	wg.Wait()
}

// 方案4. 方案2的优化，使用rand pool，保证一定并行度的同时减少内存分配
// 无并发 BenchmarkMathPool-8   	35656502	        28.68 ns/op	       0 B/op	       0 allocs/op
// 并发 BenchmarkMathPool-8   	 5006486	       237.1 ns/op	      24 B/op	       1 allocs/op
func BenchmarkMathPool(b *testing.B) {
	p := sync.Pool{
		New: func() interface{} {
			return rand.New(rand.NewSource(time.Now().UnixNano()))
		},
	}
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var tid [16]byte
			r := p.Get().(*rand.Rand)
			_, _ = r.Read(tid[:])
			p.Put(r)
		}()
	}
	wg.Wait()
}

func TestAddUintOverflow(t *testing.T) {
	var cycle uint32 = math.MaxUint32
	new := atomic.AddUint32(&cycle, 1)
	assert.Equal(t, uint32(0), new)
	new = atomic.AddUint32(&cycle, 1)
	assert.Equal(t, uint32(1), new)
}

// 总结：方案3、4其实类似，都是多个 rand 对象减少锁冲突，但是 sync.pool 可以在并发度不高的情况下回收内存节省资源

type ShardedRander struct {
	cycle   uint32
	randers []*safeRander
}

func NewShardedRander(n int) *ShardedRander {
	if n <= 0 {
		n = 1
	}
	sr := &ShardedRander{
		randers: make([]*safeRander, 0, n),
	}
	for i := 0; i < n; i++ {
		sr.randers = append(sr.randers, newSafeRander())
	}
	return sr
}

func (sc *ShardedRander) index() int {
	atomic.AddUint32(&sc.cycle, 1)
	return int(sc.cycle) % len(sc.randers)
}

func (sc *ShardedRander) Int32() int32 {
	return sc.randers[sc.index()].Int32()
}

func (sc *ShardedRander) Int32n(n int32) int32 {
	return sc.randers[sc.index()].Int32n(n)
}

func (sc *ShardedRander) Int63() int64 {
	return sc.randers[sc.index()].Int63()
}

func (sc *ShardedRander) Uint64() uint64 {
	return sc.randers[sc.index()].Uint64()
}

func (sc *ShardedRander) Read(p []byte) (n int, err error) {
	return sc.randers[sc.index()].Read(p)
}

type safeRander struct {
	item *rand.Rand
	mu   sync.Mutex
}

func newSafeRander() *safeRander {
	return &safeRander{
		item: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (sr *safeRander) Read(p []byte) (n int, err error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	return sr.item.Read(p)
}

func (sr *safeRander) Int32() int32 {
	sr.mu.Lock()
	n := sr.item.Int31()
	sr.mu.Unlock()
	return n
}

func (sr *safeRander) Int32n(n int32) int32 {
	sr.mu.Lock()
	n = sr.item.Int31n(n)
	sr.mu.Unlock()
	return n
}

func (sr *safeRander) Int63() (n int64) {
	sr.mu.Lock()
	n = sr.item.Int63()
	sr.mu.Unlock()
	return
}

func (sr *safeRander) Uint64() (n uint64) {
	sr.mu.Lock()
	n = sr.item.Uint64()
	sr.mu.Unlock()
	return
}
