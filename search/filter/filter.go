package filter

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/kelindar/bitmap"
)

type Filter32 struct {
	bm bitmap.Bitmap
	rb *roaring.Bitmap
}

// 64位的roaring bitmap, bitmap 开销太大
type Filter64 struct {
	r64 *roaring64.Bitmap
}

// type BloomFilter struct {
// 	redis *redis.Client
// }

// 初始化
func New32Filter(datas []uint32) *Filter32 {
	f := &Filter32{
		bm: bitmap.Bitmap{},
		rb: roaring.New(),
	}
	f.rb = roaring.BitmapOf(datas...)
	for _, data := range datas {
		f.bm.Set(data)
	}
	return f
}

// 初始化
func New64Filter(datas []uint64) *Filter64 {
	f := &Filter64{
		r64: roaring64.BitmapOf(datas...),
	}
	return f
}

// bitmap 过滤方法
func (f *Filter32) IsExistByBitmap(data uint32) bool {
	return f.bm.Contains(data)
}

// roaring bitmap 过滤方法
func (f *Filter32) IsExistByRoaring(data uint32) bool {
	return f.rb.Contains(data)
}

func (f *Filter64) IsExistByRoaring(data uint64) bool {
	return f.r64.Contains(data)
}

// 实现
// roaring bitmap 作为posting list
// 存储 redis:
// 方案1.对于int32的数据，可通过redis 组件实现（仅支持int32）
// 方案2.如果是int64的数据，比如雪花id，需要在内存中实现序列化到redis

func SetRedisRoaringBitmap32() {
}

func SetRedisRoaringBitmap64() {
}
