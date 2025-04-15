/**
 * @Author: wangxinyu
 * @Date: 2024/8/29 11:14
 */
package registry

import (
	"fmt"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
)

type DispatchRedis struct {
	// cfg     RedisDbDesc
	rSingle map[string][]*DbDescriptor
}

type DbDescriptor struct {
	DbName   string
	conn     *redis.Client
	DbID     int
	UserMode int
}

type ShardFunc func(key string, length int) (uint64, error)

func (t *DispatchRedis) GetDb(dbName string, key string, fun ShardFunc) (*redis.Client, error) {
	descDb, ok := t.rSingle[dbName]
	if !ok {
		return nil, fmt.Errorf("no found db with name %s", dbName)
	}
	// 分片hash算法
	rdb, err := fun(key, len(descDb))
	if err != nil {
		return nil, err
	}
	return descDb[rdb].conn, nil
}

// doConnectSingle 链接数据库
func (t *DispatchRedis) doConnectSingle(dbAlias string, cfg []*RedisConfig) *redis.Client {
	var dbs []*DbDescriptor
	for _, v := range cfg {
		address := fmt.Sprintf("%s:%d", v.Host, v.Port)
		cc := &redis.Options{
			Addr:                  address,
			Password:              v.Pwd,
			DB:                    v.Dbid,
			DialTimeout:           0,
			ReadTimeout:           3 * time.Second,
			WriteTimeout:          3 * time.Second,
			ContextTimeoutEnabled: false,
			PoolFIFO:              false,
			PoolSize:              10 * runtime.GOMAXPROCS(0),
			PoolTimeout:           0,
			MinIdleConns:          2,
			MaxIdleConns:          8,
			MaxActiveConns:        256,
		}
		cli := redis.NewClient(cc)
		dbs = append(dbs, &DbDescriptor{
			DbName:   dbAlias,
			conn:     cli,
			DbID:     v.Dbid,
			UserMode: v.Useridmod,
		})
	}
	t.rSingle[dbAlias] = dbs
	return nil
}

type RedisDbDesc struct {
	Dbid      int    `json:"dbid"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Pwd       string `json:"pwd"`
	UserIdMod int    `json:"useridmod"`
}

var (
	rc *DispatchRedis
)

// MusConnectRedis connect redis server
func MusConnectRedis(cfg map[string][]*RedisConfig) {
	tr := new(DispatchRedis)
	tr.rSingle = map[string][]*DbDescriptor{}
	for k, v := range cfg {
		tr.doConnectSingle(k, v)
	}
	rc = tr
}
