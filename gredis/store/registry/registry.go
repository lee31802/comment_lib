package registry

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"hash/fnv"
	"math"
	"strconv"
	"strings"
)

func (hr *HashRedisClient) hash(data string) (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(data))
	if err != nil {
		// TODO 使用errors
		return 0, err
	}
	return h.Sum32(), nil
}

func (hr *HashRedisClient) GetRedisClient(key string) (*redis.Client, error) {
	return rc.GetDb(hr.name, key, func(key string, length int) (uint64, error) {
		// 多个库的情况下需要按照key 进行hash取mod
		hashKey, err := hr.hash(key)
		if err != nil {
			return 999, err
		}
		// TODO 32位, 64各种位转换看起来有点乱
		idx := uint64(math.Abs(float64(hashKey)))
		return idx % uint64(length), nil
	})
}

type SingleRedisClient struct {
	name string //name 私有化防止运行期误操作
}

func (sr *SingleRedisClient) GetRedisClient(key string) (*redis.Client, error) {
	return rc.GetDb(sr.name, key, func(key string, length int) (uint64, error) {
		if length == 1 {
			return 0, nil
		}
		// TODO 使用errors
		return 999, fmt.Errorf("singleRedis error %s %d", key, length)
	})
}

// HashRedisClient 按照key进行hash确定db, 兼容单DB
type HashRedisClient struct {
	// TODO key 目前看都没有用
	name string
	Key  string
}

// UserIdRedisClient 按照userId取模确定db
type UserIdRedisClient struct {
	name string
}

func (ur *UserIdRedisClient) parserUidByKey(key string) (uint64, error) {
	strList := strings.Split(key, ":")
	if len(strList) != 2 {
		return 0, fmt.Errorf("invalid user redis key: %s", key)
	}
	return strconv.ParseUint(strList[1], 10, 64)
}

func (ur *UserIdRedisClient) GetRedisClient(key string) (*redis.Client, error) {
	return rc.GetDb(ur.name, key, func(key string, length int) (uint64, error) {
		userId, err := ur.parserUidByKey(key)
		if err != nil {
			return 999, err
		}
		return userId % uint64(length), nil
	})
}

func NewHashRedisClient(name string) *HashRedisClient {
	return &HashRedisClient{
		name: name,
	}
}

func NewUserIdRedisClient(name string) *UserIdRedisClient {
	return &UserIdRedisClient{
		name: name,
	}
}

func NewSingleRedisClient(name string) *SingleRedisClient {
	return &SingleRedisClient{
		name: name,
	}
}

type Register interface {
	// GetRedisClient TODO name 参数有点多余， 实现类上已经有名称了
	GetRedisClient(key string) (*redis.Client, error)
}
