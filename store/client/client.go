package client

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golog/comm_lib/store/registry"
	"reflect"
	"time"
)

var Nil = redis.Nil

type StoreClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	dbs    registry.Register
}

func (s *StoreClient) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	// TODO 这个保护感觉必要性不大
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	// TODO 空指针异常
	return cli.ZAdd(ctx, key, members...)
}

func (s *StoreClient) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.ZCount(ctx, key, min, max)
}

func (s *StoreClient) ZCard(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.ZCard(ctx, key)
}

func (s *StoreClient) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.ZRem(ctx, key, members...)
}

func (s *StoreClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.TTL(ctx, key)
}

func (s *StoreClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.Incr(ctx, key)
}

func (s *StoreClient) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.IncrBy(ctx, key, value)
}

func (s *StoreClient) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.IncrByFloat(ctx, key, value)
}

func (s *StoreClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.Decr(ctx, key)
}

func (s *StoreClient) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.DecrBy(ctx, key, decrement)
}

func (s *StoreClient) Append(ctx context.Context, key, value string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.Append(ctx, key, value)
}

func (s *StoreClient) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
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

func (s *StoreClient) HMSet(ctx context.Context, key string, resp interface{}) *redis.BoolCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	fields := GetFieldNamesAndValues(resp)
	return cli.HMSet(ctx, key, fields...)
}

func (s *StoreClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.HSet(ctx, key, values...)
}

func (s *StoreClient) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.LPush(ctx, key, values...)
}

func (s *StoreClient) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.LRange(ctx, key, start, stop)
}

// Set string set命令
// expires 单位为秒
func (s *StoreClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.Set(ctx, key, value, expiration)
}

func (s *StoreClient) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.GetSet(ctx, key, value)
}

func (s *StoreClient) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.SetEx(ctx, key, value, expiration)
}

func (s *StoreClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.SetNX(ctx, key, value, expiration)
}

func (s *StoreClient) toBoolean(val string) bool {
	r := val
	return r == "OK" || r == "1" || r == "true"
}

func (s *StoreClient) Get(ctx context.Context, key string) *redis.StringCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.Get(ctx, key)
}
func (s *StoreClient) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.HGetAll(ctx, key)
}

func (s *StoreClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
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

func (s *StoreClient) HMGet(ctx context.Context, key string, resp interface{}) *redis.SliceCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	fields := GetRedisTags(resp)
	//glog.XDebugf("GetRedisTags result %v, %v", key, fields)
	//rst, err := cli.HMGet(ctx, key, fields...).Result()
	//glog.XDebugf("GetRedisTags redis result %v, %v %v, %v", key, fields, rst, err)
	return cli.HMGet(ctx, key, fields...)
}

func (s *StoreClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.SAdd(ctx, key, members...)
}

func (s *StoreClient) SMIsMember(ctx context.Context, key string, members ...interface{}) *redis.BoolSliceCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.SMIsMember(ctx, key, members...)
}

func (s *StoreClient) Del(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		//	todo
	}
	return cli.Del(ctx, key)
}

// Exists 判断key是否存在
func (s *StoreClient) Exists(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		// todo
	}
	return cli.Exists(ctx, key)
}

func (s *StoreClient) Expire(ctx context.Context, key string, duration time.Duration) *redis.BoolCmd {
	if ctx == nil {
		ctx = context.Background()
	}
	cli, err := s.dbs.GetRedisClient(key)
	if cli == nil || err != nil {
		// todo
	}
	return cli.Expire(ctx, key, duration)
}

func NewHashClient(dbName string) *StoreClient {
	sc := new(StoreClient)
	// TODO 连接上做5s超时干嘛？
	sc.ctx, sc.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	sc.dbs = registry.NewHashRedisClient(dbName)
	return sc
}
func NewSingleClient(dbName string) *StoreClient {
	sc := new(StoreClient)
	sc.ctx, sc.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	sc.dbs = registry.NewSingleRedisClient(dbName)
	return sc
}
func NewUserClient(dbName string) *StoreClient {
	sc := new(StoreClient)
	sc.ctx, sc.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	sc.dbs = registry.NewUserIdRedisClient(dbName)
	return sc
}
