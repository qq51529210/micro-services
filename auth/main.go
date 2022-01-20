package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/qq51529210/log"
	"github.com/qq51529210/micro-services/auth/cache"
	"github.com/qq51529210/micro-services/auth/cache/redis"
	"github.com/qq51529210/micro-services/auth/service/api"
	"github.com/qq51529210/micro-services/auth/service/oauth2"
	"github.com/qq51529210/micro-services/auth/store"
	"github.com/qq51529210/micro-services/auth/store/mongodb"
	"github.com/qq51529210/web"
	"github.com/qq51529210/web/router"
)

var (
	caches = make(map[string]func(map[string]interface{}) cache.Cache)
	stores = make(map[string]func(map[string]interface{}) store.Store)
)

type config struct {
	Listen      string        `json:"listen"`
	X509CertPEM []string      `json:"x509CertPEM"`
	X509KeyPEM  []string      `json:"x509KeyPEM"`
	RootDir     string        `json:"rootDir"`
	Cookie      *cookieConfig `json:"cookie"`
	UUID        *uuidConfig   `json:"uuid"`
	Store       *storeConfig  `json:"store"`
	Cache       *storeConfig  `json:"cache"`
}

type cookieConfig struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
	MaxAge int64  `json:"maxAge"`
}

type uuidConfig struct {
	GroupID   byte   `json:"groupID"`
	MechineID byte   `json:"mechineID"`
	Node      string `json:"node"`
}

type storeConfig struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

func init() {
	caches[""] = redis.New
	caches["redis"] = redis.New
	//
	stores[""] = mongodb.New
	stores["mongodb"] = mongodb.New
}

func loadConfig() *config {
	var path string
	if len(os.Args) < 2 {
		// 如果没有参数，则尝试加载当前目录下的名为"{app-name}.json"的文件。
		dir, file := filepath.Split(os.Args[0])
		ext := filepath.Ext(file)
		if ext != "" {
			file = file[:len(file)-len(ext)]
		}
		path = filepath.Join(dir, file+".json")
	} else {
		// 使用命令行第一个参数作为路径(本地磁盘/http资源)。
		path = os.Args[1]
	}
	var data []byte
	var err error
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		// http路径
		res, err := http.Get(path)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
	} else {
		// 本地路径
		data, err = ioutil.ReadFile(path)
	}
	if err != nil {
		panic(err)
	}
	// 解析json
	cf := new(config)
	err = json.Unmarshal(data, cf)
	if err != nil {
		panic(err)
	}
	return cf
}

func initRouter(cfg *config) router.RootRouter {
	root := router.NewRootRouter()
	// 静态
	root.Static("static", cfg.RootDir, 1024*1024)
	//
	oauth2.Init(root.SubRouter("oauth2"))
	api.Init(root.SubRouter("api"))
	//
	return root
}

func initServer() web.Server {
	// 加载启动配置
	conf := loadConfig()
	// 缓存
	newCache := caches[conf.Cache.Type]
	if newCache == nil {
		panic(fmt.Errorf("config.cache.type: unsupported cache <%s>", conf.Cache.Type))
	}
	cache.SetCache(newCache(conf.Cache.Config))
	// 数据库
	newStore := caches[conf.Store.Type]
	if newStore == nil {
		panic(fmt.Errorf("config.cache.type: unsupported store <%s>", conf.Store.Type))
	}
	store.SetStore(newStore(conf.Store.Config))
	// http路由和handler
	router := initRouter(conf)
	// 服务
	var ser web.Server
	if len(conf.X509CertPEM) != 0 && len(conf.X509KeyPEM) != 0 {
		ser = web.NewTLSServerWithKeyPair(conf.Listen,
			[]byte(strings.Join(conf.X509CertPEM, "")),
			[]byte(strings.Join(conf.X509KeyPEM, "")), router)
	} else {
		ser = web.NewServer(conf.Listen, router)
	}
	return ser
}

func main() {
	defer func() {
		re := recover()
		if re != nil {
			log.Recover(re)
		}
	}()
	ser := initServer()
	err := ser.Serve()
	if err != nil {
		log.Error(err)
	}
}
