package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/go-leo/sonyflake"
	"github.com/go-redis/redis/v8"
)

// 1.雪花算法ID bitmap实践
// 2.roraingbitmap redis实践
// 3.boolmfilter
// 4.延迟队列
// 5.hoylog
// 6.地理

var (
	rdb *redis.Client

	sf2014 *sonyflake.Sonyflake
	sf2015 *sonyflake.Sonyflake
	sf2023 *sonyflake.Sonyflake
)

// Init redis conn
func init() {
	initSonyflake()
	// 创建一个 Redis 客户端
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器的地址和端口
		Password: "",               // 如果没有密码，使用空字符串
		DB:       0,                // 使用默认的 DB
	})
}

func initSonyflake() {
	rand.Seed(time.Now().Unix())
	st := sonyflake.Settings{
		StartTime: time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC),
		Sequence:  func() uint16 { return uint16(rand.Intn(1024)) },
	}
	sf2014 = sonyflake.NewSonyflake(st)
	if sf2014 == nil {
		panic("sonyflake not created")
	}
	st1 := sonyflake.Settings{
		StartTime: time.Date(2015, 9, 1, 0, 0, 0, 0, time.UTC),
		Sequence:  func() uint16 { return uint16(rand.Intn(1024)) },
	}
	sf2015 = sonyflake.NewSonyflake(st1)
	if sf2015 == nil {
		panic("sonyflake not created")
	}
	st2 := sonyflake.Settings{
		StartTime: time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
		Sequence:  func() uint16 { return uint16(rand.Intn(1024)) },
	}
	sf2023 = sonyflake.NewSonyflake(st2)
	if sf2015 == nil {
		panic("sonyflake not created")
	}
}

func main() {
	// UseBitmap()
	UseRoraingBitmap()
}

// bitmap
// 从系统运行起的第一个id，确定为起始范围
// 一亿一个bitmap范围 100000000 / 8 / 1024 / 1024 = 10mb
// 不行，每过一秒就有10亿个id过去，太过于稀疏（一个id就占用10mb）
func UseBitmap() {
	for i := 0; i < 100; i++ {
		n := time.Duration(rand.Intn(3)) * time.Second
		time.Sleep(n)
		id, _ := sf2014.NextID()
		shard := (id - 483373146767901013) / 100000000
		offset := int64((id - 483373146767901013) % 100000000)
		rdb.SetBit(context.Background(), fmt.Sprintf("bigid:%d", shard), offset, 1)
	}
} // 42  270mb

// 方案一：尝试roraing bitmap + 分片
// 按照uid范围分片，序列化后存储
// tostring tobase64 和普通字符串存储大小没有太大区别 {483947267545383857,483947267545383858,483947267545383859,483947267545383860,483947269223104676}
// 只是在内存中的大小不同，但是要存储在redis中只能序列话后存储（不支持roraingbitmap）内存中 748 B， redis中 1440 B
func UseRoraingBitmap() {
	r64 := roaring64.NewBitmap()
	shard := 0
	for i := 0; i < 100; i++ {
		n := time.Duration(rand.Intn(3)) * time.Second
		time.Sleep(n)
		id, _ := sf2014.NextID()
		fmt.Println(id)
		r64.Add(id)
		// if r64.GetSerializedSizeInBytes() > 1024*1024*10 {
		// 	fmt.Printf("roraing bitmap too large: %d, size: %d, static:%+v\n", r64.GetSerializedSizeInBytes(), r64.GetSizeInBytes(), r64.Stats())
		// 	rdb.Set(context.Background(), fmt.Sprintf("bigid:%d", shard), r64.String(), 0)
		// 	shard++
		// 	r64.Clear()
		// }
	}
	if r64.GetSerializedSizeInBytes() > 0 {
		fmt.Printf("end roraing bitmap too large: %d, size: %d, static:%+v\n", r64.GetSerializedSizeInBytes(), r64.GetSizeInBytes(), r64.Stats())
		ret, _ := r64.ToBase64()
		rdb.Set(context.Background(), fmt.Sprintf("bigid:%d", shard), ret, 0)
	}
}

// 优化雪花算法？，减少时间、设备、随机数的长度；没有解决根本问题，id不连续，上限还是高；

// 方案二：提供一个发号器服务，每次生成一个递增id，并且给老的id都生成一个短id（类似于短url）
