package registry

import (
	"fmt"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
)

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

// NewClient 根据配置创建对应的Redis客户端
func NewClient(cfg Config) (RedisClient, error) {
	switch cfg.Mode {
	case "cluster":
		return newClusterClient(cfg.Cluster)
	case "ring":
		return newRingClient(cfg.Ring)
	case "sentinel":
		return newFailoverClient(cfg.Sentinel)
	case "single":
		return newDefaultClient(cfg.Standalone)
	default:
		return nil, fmt.Errorf("unsupported redis mode: %s", cfg.Mode)
	}
}
