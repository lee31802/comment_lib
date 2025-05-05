package env

//
//import (
//	"fmt"
//	"net"
//	"os"
//	"os/exec"
//	"strings"
//
//	"github.com/lee31802/comment_lib/logkit"
//	"go.uber.org/zap"
//)
//
//const (
//	defaultCID   = "local"
//	defaultEnv   = "dev"
//	defaultIndex = "0"
//	defaultZKURL = "zookeeper-1:2181,zookeeper-2:2181,zookeeper-3:2181"
//)
//
//var (
//	globalEnv Env
//)
//
//func init() {
//	globalEnv.LoadConfig()
//}
//
//type Env struct {
//	CID      string `json:"cid"`
//	Env      string `json:"env"`
//	Host     string `json:"host"`
//	Port     string `json:"port"`
//	APPID    string `json:"appid"`
//	Project  string `json:"project"`
//	Module   string `json:"module"`
//	Service  string `json:"service"`
//	Index    string `json:"index"`
//	ZKURL    string `json:"zkurl"`
//	HostName string `json:"host_name"`
//	PfbName  string `json:"pfb_name"`
//	IsLive   bool   `json:"is_live"`
//}
//
//func (env *Env) ReadEnvLower(key string) string {
//	return strings.TrimSpace(strings.ToLower(os.Getenv(key)))
//}
//
//func (env *Env) LoadConfig() {
//	env.CID = env.ReadEnvLower("cid")
//	env.Env = env.ReadEnvLower("env")
//	env.Host = env.ReadEnvLower("HOST")
//	env.Port = env.ReadEnvLower("PORT")
//	env.APPID = env.ReadEnvLower("MARATHON_APP_ID")
//	env.Project = env.ReadEnvLower("PROJECT_NAME")
//	env.Module = env.ReadEnvLower("MODULE_NAME")
//	env.Service = fmt.Sprintf("%s-%s", env.Project, env.Module)
//	env.Index = env.ReadEnvLower("INDEX")
//	env.ZKURL = env.ReadEnvLower("ZK_URL")
//	env.HostName = env.ReadEnvLower("HOSTNAME")
//	env.PfbName = env.ReadEnvLower("PFB_NAME")
//
//	if env.CID == "" {
//		env.CID = defaultCID
//		logkit.Warn("CID no found in env, set default:", zap.String("CID", env.CID))
//	}
//	if env.Env == "" {
//		env.Env = defaultEnv
//		logkit.Warn("Env no found in env, set default:", zap.String("Env", env.Env))
//	}
//	if env.Host == "" {
//		env.Host = env.getIP()
//		logkit.Warn("Host no found in env, set default:", zap.String("Host", env.Host))
//	}
//	if env.Index == "" {
//		env.Index = defaultIndex
//	}
//	if env.ZKURL == "" {
//		env.ZKURL = defaultZKURL
//	}
//	env.IsLive = env.Env == "live"
//}
//
//func (env *Env) getIP() string {
//	if env.Env == defaultEnv {
//		logkit.Warn("getLocalIP", logkit.String("env", env.Env))
//		return env.getLocalIP()
//	}
//	return env.getDockerIP()
//}
//
//func (env *Env) getLocalIP() string {
//	addrs, err := net.InterfaceAddrs()
//
//	if err != nil {
//		logkit.Error("get local ip failed", zap.Error(err))
//		return ""
//	}
//
//	for _, address := range addrs {
//
//		// 检查ip地址判断是否回环地址
//		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//			if ipnet.IP.To4() != nil {
//				return ipnet.IP.String()
//			}
//		}
//	}
//	return ""
//}
//
//func (env *Env) getDockerIP() string {
//	cmd := "getent hosts internal-dns  | awk '{ print $1 }' | xargs ip -4 route get | grep dev"
//	out, err := exec.Command("bash", "-c", cmd).Output()
//	if err != nil {
//		logkit.Error("get local ip failed", logkit.Err(err))
//		return env.getLocalIP()
//	}
//	args := strings.Split(strings.TrimSpace(string(out)), " ")
//	if len(args) == 0 {
//		return env.getLocalIP()
//	}
//	return args[len(args)-1]
//}
//
//func GetCID() string {
//	return globalEnv.CID
//}
//
//func GetEnv() string {
//	return globalEnv.Env
//}
//
//func GetHost() string {
//	return globalEnv.Host
//}
//
//func GetPort() string {
//	return globalEnv.Port
//}
//
//func GetProject() string {
//	return globalEnv.Project
//}
//
//func GetModule() string {
//	return globalEnv.Module
//}
//
//func GetService() string {
//	return globalEnv.Service
//}
//
//func GetIndex() string {
//	return globalEnv.Index
//}
//
//func GetZKURL() string {
//	return globalEnv.ZKURL
//}
//
//func GetHostName() string {
//	return globalEnv.HostName
//}
//
//func IsLive() bool {
//	return globalEnv.IsLive
//}
//func PfbName() string {
//	return globalEnv.PfbName
//}
//
//func SupportPfb() bool {
//	return !IsLive()
//}
//
//func SetGlobalEnv(env *Env) {
//	globalEnv = *env
//}
