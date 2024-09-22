/* redis数据库的相关操作，即初始化，以及缓存排行榜 */

package common

import (
	"context"
	"fmt"
	"ginEssential/controller/model"

	"github.com/redis/go-redis/v9"
)

var (
	CTX = context.Background()
	// RANK=false
	REDIS_DB *redis.Client // Redis的数据库连接
	// a redis.Z
)

func init() {
	REDIS_DB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := REDIS_DB.Set(CTX, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
}

// 音乐的分数
func MusicScore(m *model.Music) redis.Z {
	return redis.Z{
		Score:  float64(m.Poll),
		Member: m.ID,
	}
}

// 官网上的测试示例
func ConnectRedis_demo() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(CTX, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(CTX, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("The value of key is", val)

	val2, err := rdb.Get(CTX, "key2").Result()
	if err == redis.Nil { // 用来判断键是否存在，严格来说不算错误
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
