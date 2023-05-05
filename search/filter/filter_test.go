package filter

import (
	"math"
	"math/rand"
	"testing"

	"github.com/RoaringBitmap/roaring"
	"github.com/kelindar/bitmap"
)

// 一定范围连续数据的过滤
// BenchmarkFilteBitmap-8   	1000000000	         0.9560 ns/op	       0 B/op	       0 allocs/op
// 3.27MB
func BenchmarkFilteBitmap(b *testing.B) {
	datas := make([]uint32, 1000000)
	for i := 0; i < 1000000; i++ {
		datas[i] = uint32(i)
	}

	b.ResetTimer()
	var bm bitmap.Bitmap
	for _, data := range datas {
		bm.Set(data)
	}
	for i := 0; i < b.N; i++ {
		_ = bm.Contains(uint32(i))
	}
}

// 一定范围连续数据的过滤
// BenchmarkFilteRoaring-8   	70526587	        17.91 ns/op	       0 B/op	       0 allocs/op
//  3.52MB
func BenchmarkFilteRoaring(b *testing.B) {
	datas := make([]uint32, 1000000)
	for i := 0; i < 1000000; i++ {
		datas[i] = uint32(i)
	}
	b.ResetTimer()
	r := roaring.BitmapOf(datas...)
	for i := 0; i < b.N; i++ {
		_ = r.Contains(uint32(i))
	}
}

// 一定范围的随机数据的过滤
// BenchmarkFilteBitmapRandom-8    1000000000               1.058 ns/op           0 B/op          0 allocs/op
// 6.41GB
func BenchmarkFilteBitmapRandom(b *testing.B) {
	datas := make([]uint32, 1000000)
	for i := 0; i < 1000000; i++ {
		datas[i] = uint32(rand.Intn(math.MaxInt32))
	}
	b.ResetTimer()
	var bm bitmap.Bitmap
	for _, data := range datas {
		bm.Set(data)
	}
	for i := 0; i < b.N; i++ {
		_ = bm.Contains(uint32(i))
	}
}

// 一定范围的随机数据的过滤
// BenchmarkFilteRoaringRandom-8           23812920                48.55 ns/op            0 B/op          0 allocs/op
// 161.58MB
func BenchmarkFilteRoaringRandom(b *testing.B) {
	datas := make([]uint32, 1000000)
	for i := 0; i < 1000000; i++ {
		datas[i] = uint32(rand.Intn(math.MaxInt32))
	}
	b.ResetTimer()
	r := roaring.BitmapOf(datas...)
	for i := 0; i < b.N; i++ {
		_ = r.Contains(uint32(i))
	}
	// b.Run("bitmap", func(b *testing.B) {
	// 	var bm bitmap.Bitmap
	// 	for _, data := range datas {
	// 		bm.Set(data)
	// 	}
	// 	for i := 0; i < b.N; i++ {
	// 		_ = r.Contains(uint32(i))
	// 	}
	// })
}
