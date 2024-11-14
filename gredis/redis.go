package gredis

import (
    "context"
    "encoding/json"
    "gin-example/gin-blog/setting"
    "github.com/redis/go-redis/v9"
    "log"
    "strings"
    "time"
)

var rdb *redis.Client
var ctx = context.Background()

// Setup 连接redis数据库
func Setup() {
    rdb = redis.NewClient(&redis.Options{
        Addr:            setting.RedisSetting.Host,
        Password:        setting.RedisSetting.Password,
        DB:              setting.RedisSetting.DB,
        MaxIdleConns:    setting.RedisSetting.MaxIdle,
        MaxActiveConns:  setting.RedisSetting.MaxActive,
        ConnMaxIdleTime: setting.RedisSetting.IdleTimeout,
    })

    // 测试连接
    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("redis ping err: %v", err)
    }
}

func Set(key string, data interface{}, expiration time.Duration) error {
    val, err := json.Marshal(data)
    if err != nil {
        return err
    }
    err = rdb.Set(ctx, key, string(val), expiration).Err()
    return err
}

func Exists(key string) bool {
    exists, err := rdb.Exists(ctx, key).Result()
    if err != nil || exists == 0 {
        return false
    }
    return true
}

func LikeExists(key ...string) bool {
    keys := strings.Join(key, "*")
    keys = "*" + keys + "*"
    keyList, err := rdb.Keys(ctx, keys).Result()
    if err != nil {
        return false
    }
    return len(keyList) > 0
}

func Get(key string) (string, error) {
    val, err := rdb.Get(ctx, key).Result()
    return val, err
}

func Del(key string) error {
    return rdb.Del(ctx, key).Err()
}

func LikeDel(key ...string) error {
    keys := strings.Join(key, "*")
    keys = "*" + keys + "*"
    keyList, err := rdb.Keys(ctx, keys).Result()
    if err != nil {
        return err
    }
    for _, k := range keyList {
        err = rdb.Del(ctx, k).Err()
        if err != nil {
            return err
        }
    }
    return nil
}
