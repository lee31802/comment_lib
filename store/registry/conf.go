/**
 * @Author: wangxinyu
 * @Date: 2024/9/30 15:55
 */
package registry

type RedisConfig struct {
	Dbid      int    `json:"dbid"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Pwd       string `json:"pwd"`
	Useridmod int    `json:"useridmod"`
}
