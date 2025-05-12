package client

import (
	"context"
	"fmt"
	"reflect"
	"sdk-go/pkg/store/registry"
	"time"

	"github.com/redis/go-redis/v9"
)

var Nil = redis.Nil

type RedisClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	dbs    registry.Register
}

func (s *RedisClient) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	// TODO 这个保护感觉必要性不大
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	// TODO 空指针异常
	return cli.ZAdd(ctx, key, members...)
}

func (s *RedisClient) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.ZCount(ctx, key, min, max)
}

func (s *RedisClient) ZCard(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.ZCard(ctx, key)
}

func (s *RedisClient) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.ZRem(ctx, key, members...)
}

func (s *RedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.TTL(ctx, key)
}

func (s *RedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.Incr(ctx, key)
}

func (s *RedisClient) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.IncrBy(ctx, key, value)
}

func (s *RedisClient) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.IncrByFloat(ctx, key, value)
}

func (s *RedisClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.Decr(ctx, key)
}

func (s *RedisClient) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.DecrBy(ctx, key, decrement)
}

func (s *RedisClient) Append(ctx context.Context, key, value string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.Append(ctx, key, value)
}

func (s *RedisClient) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.HKeys(ctx, key)
}

func GetFieldNamesAndValues(t interface{}) []interface{} {
	tValue := reflect.ValueOf(t).Elem()
	tType := tValue.Type()

	var result []interface{}
	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		value := tValue.Field(i)

		// Get the redis tag or fallback to field name
		fieldName := field.Tag.Get("redis")
		if fieldName == "" {
			fieldName = field.Name
		}

		// Append field name and value to the result slice
		result = append(result, fieldName, value.Interface())
	}

	return result
}

func (s *RedisClient) HMSet(ctx context.Context, key string, resp interface{}) *redis.BoolCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	fields := GetFieldNamesAndValues(resp)
	return cli.HMSet(ctx, key, fields...)
}

func (s *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.HSet(ctx, key, values...)
}

func (s *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.LPush(ctx, key, values...)
}

func (s *RedisClient) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.LRange(ctx, key, start, stop)
}

// Set string set命令
// expires 单位为秒
func (s *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.Set(ctx, key, value, expiration)
}

func (s *RedisClient) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.GetSet(ctx, key, value)
}

func (s *RedisClient) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.SetEx(ctx, key, value, expiration)
}

func (s *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.SetNX(ctx, key, value, expiration)
}

func (s *RedisClient) toBoolean(val string) bool {
	r := val
	return r == "OK" || r == "1" || r == "true"
}

func (s *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.Get(ctx, key)
}
func (s *RedisClient) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.HGetAll(ctx, key)
}

func (s *RedisClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.HGet(ctx, key, field)
}

func GetRedisTags(t interface{}) []string {
	if t == nil {
		return []string{}
	}
	tType := reflect.TypeOf(t)

	// Check if the input is a pointer and get the element type.
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}

	if tType.Kind() != reflect.Struct {
		fmt.Println("Input is not a struct type")
		return nil
	}

	var tags []string
	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		redisTag := field.Tag.Get("redis")
		if redisTag != "" {
			tags = append(tags, redisTag)
		}
	}
	return tags
}

func (s *RedisClient) HMGet(ctx context.Context, key string, resp interface{}) *redis.SliceCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	fields := GetRedisTags(resp)
	//glog.XDebugf("GetRedisTags result %v, %v", key, fields)
	//rst, err := cli.HMGet(ctx, key, fields...).Result()
	//glog.XDebugf("GetRedisTags redis result %v, %v %v, %v", key, fields, rst, err)
	return cli.HMGet(ctx, key, fields...)
}

func (s *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.SAdd(ctx, key, members...)
}

func (s *RedisClient) SMIsMember(ctx context.Context, key string, members ...interface{}) *redis.BoolSliceCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.SMIsMember(ctx, key, members...)
}

func (s *RedisClient) Del(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.Del(ctx, key)
}

// Exists 判断key是否存在
func (s *RedisClient) Exists(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.Exists(ctx, key)
}

func (s *RedisClient) Expire(ctx context.Context, key string, duration time.Duration) *redis.BoolCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		glog.XInfo("get redis client is nil err : %v", err)
	}
	return cli.Expire(ctx, key, duration)
}

func NewRingClient(dbName string) *RedisClient {
	sc := new(RedisClient)
	// TODO 连接上做5s超时干嘛？
	sc.ctx, sc.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	sc.dbs = registry.NewHashRedisClient(dbName)
	return sc
}
func NewClusterClient(dbName string) *RedisClient {
	sc := new(RedisClient)
	sc.ctx, sc.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	sc.dbs = registry.NewSingleRedisClient(dbName)
	return sc
}
func NewFailoverClient(dbName string) *RedisClient {
	sc := new(RedisClient)
	sc.ctx, sc.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	sc.dbs = registry.NewUserIdRedisClient(dbName)
	return sc
}

func NewDefaultClient(dbName string) *RedisClient {
	sc := new(RedisClient)
	sc.ctx, sc.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	sc.dbs = registry.NewUserIdRedisClient(dbName)
	return sc
}
