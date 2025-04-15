/**
 * @Author: wangxinyu
 * @Date: 2024/9/18 18:30
 */
package store

import (
	"golog/comm_lib/store/client"
)

const (
	MixDb       = "mix"
	UserDb      = "user"
	ConfigDb    = "config"
	AvatarDb    = "avatar"
	PayDataDb   = "paydata"
	UserKeysDb  = "userkeys"
	BiCountDb   = "bicount"
	friendDb    = "friend"
	lockerDb    = "locker"
	forbiddenDb = "forbidden"
	SdkDataDb   = "sdkdatas"
)

func Mix() *client.StoreClient {
	cli := client.NewSingleClient(MixDb)
	return cli
}

func User() *client.StoreClient {
	return client.NewUserClient(UserDb)
}

func Config() *client.StoreClient {
	cli := client.NewSingleClient(ConfigDb)
	return cli
}

func Avatar() *client.StoreClient {
	cli := client.NewSingleClient(AvatarDb)
	return cli
}

func PayData() *client.StoreClient {
	cli := client.NewSingleClient(PayDataDb)
	return cli
}

func UserKeys() *client.StoreClient {
	return client.NewHashClient(UserKeysDb)
}

func BiCount() *client.StoreClient {
	cli := client.NewSingleClient(BiCountDb)
	return cli
}

func FriendMix() *client.StoreClient {
	cli := client.NewSingleClient(friendDb)
	return cli
}

func Locker() *client.StoreClient {
	cli := client.NewSingleClient(lockerDb)
	return cli
}

func Forbidden() *client.StoreClient {
	cli := client.NewSingleClient(forbiddenDb)
	return cli
}

func SdkDatas() *client.StoreClient {
	return client.NewHashClient(SdkDataDb)
}

//func getDbDesc() map[string][]*registry.RedisConfig {
//	return map[string][]*registry.RedisConfig{
//		AvatarDb:    {copyConfig(globalconf.TyGlobal().RedisAvatar())},
//		UserDb:      copySliceConfig(globalconf.TyGlobal().RedisDatas()),
//		MixDb:       {copyConfig(globalconf.TyGlobal().RedisMix())},
//		ConfigDb:    {copyConfig(globalconf.TyGlobal().RedisConfig())},
//		PayDataDb:   {copyConfig(globalconf.TyGlobal().RedisPaydata())},
//		UserKeysDb:  copySliceConfig(globalconf.TyGlobal().RedisUserkeys()),
//		BiCountDb:   {copyConfig(globalconf.TyGlobal().RedisBicount())},
//		friendDb:    {copyConfig(globalconf.TyGlobal().RedisFriend())},
//		lockerDb:    {copyConfig(globalconf.TyGlobal().RedisLocker())},
//		forbiddenDb: {copyConfig(globalconf.TyGlobal().RedisForbidden())},
//		SdkDataDb:   copySliceConfig(globalconf.TyGlobal().RedisSdkdatas()),
//	}
//}
//
//func copyConfig(config *globalconf.RedisConfig) *registry.RedisConfig {
//	return &registry.RedisConfig{
//		Dbid:      config.Dbid,
//		Host:      config.Host,
//		Port:      config.Port,
//		Pwd:       config.Pwd,
//		Useridmod: config.Useridmod,
//	}
//}
//func copySliceConfig(config []*globalconf.RedisConfig) []*registry.RedisConfig {
//	var cs []*registry.RedisConfig
//	for _, v := range config {
//		cs = append(cs, copyConfig(v))
//	}
//	return cs
//}
//
//func InitRedis() {
//	registry.MusConnectRedis(getDbDesc())
//}
