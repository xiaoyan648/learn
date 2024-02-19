package main

// 官方文档 https://github.com/dtm-labs/rockscache/blob/main/helper/README-cn.md
// 首个确保最终一致、强一致的 Redis 缓存库。

import (
	"log"

	"github.com/dtm-labs/rockscache"
)

// var rocksKey = "client-test-key"

func RocksGet() {
	rockscache.SetVerbose(true)
	// 使用默认选项生成rockscache的客户端
	rc := rockscache.NewClient(rdb, rockscache.NewDefaultOptions())

	// 使用Fetch获取数据，第一个参数是数据的key，第二个参数为数据过期时间，第三个参数为缓存不存在时，数据获取函数
	v, err := rc.Fetch("key1", 300, func() (string, error) {
		// 从数据库或其他渠道获取数据
		return "value1", nil
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(v)

	if err := rc.TagAsDeleted("key1"); err != nil {
		log.Fatal(err)
		return
	}
}
