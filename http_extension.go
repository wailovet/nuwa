package nuwa

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// 配置文件

type UniversalConfigType string

// sqlite
const UniversalConfigTypeSqlite UniversalConfigType = "sqlite"

// bolt
const UniversalConfigTypeBolt UniversalConfigType = "bolt"

type NuwaUniversalConfig struct {
	Id    int `xorm:"pk autoincr"`
	Key   string
	Value string
}

type NuwaUniversalConfigInterface struct {
	_get func(key string) string
	_set func(key string, value string)
}

func (c *NuwaUniversalConfigInterface) Get(key string) string {
	return c._get(key)
}
func (c *NuwaUniversalConfigInterface) Set(key string, value string) {
	c._set(key, value)
}

func (h *HttpEngine) EnableUniversalConfig(configType UniversalConfigType) *NuwaUniversalConfigInterface {

	if configType == UniversalConfigTypeSqlite {
		Sqlited().Xorm().Sync2(&NuwaUniversalConfig{})
	}

	configer := NuwaUniversalConfigInterface{
		_get: func(key string) string {
			if configType == UniversalConfigTypeSqlite {
				data := NuwaUniversalConfig{}
				Sqlited().Xorm().Where("key = ?", key).Get(&data)
				return data.Value
			} else if configType == UniversalConfigTypeBolt {
				return Bolt().Bucket("nuwa_universal_config").GetRaw(key)
			}
			return ""
		},
		_set: func(key string, value string) {
			if configType == UniversalConfigTypeSqlite {
				data := NuwaUniversalConfig{}
				Sqlited().Xorm().Where("key = ?", key).Get(&data)
				if data.Id > 0 {
					data.Value = value
					Sqlited().Xorm().Update(&data)
				} else {
					data.Key = key
					data.Value = value
					Sqlited().Xorm().Insert(&data)
				}
			} else if configType == UniversalConfigTypeBolt {
				Bolt().Bucket("nuwa_universal_config").SetRaw(key, value)
			}
		},
	}
	//
	h.HandleFunc("/nuwa/universal/config/get", func(ctx HttpContext) {
		key := ctx.ParamRequired("key")
		value := configer.Get(key)
		ctx.DisplayByData(value)
	})

	h.HandleFunc("/nuwa/universal/config/set", func(ctx HttpContext) {
		key := ctx.ParamRequired("key")
		value := ctx.ParamRequired("value")
		configer.Set(key, value)
		ctx.DisplayBySuccess()
	})

	return &configer
}

// 代理转发
func reverseProxy(addr string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	origin := ""
	isOptions := false
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		origin = req.Header.Get("Origin")
		if req.Method == "OPTIONS" {
			isOptions = true
		}
	}

	proxy.ModifyResponse = func(r *http.Response) error {

		if isOptions {
			r.StatusCode = 200
		}

		r.Header.Set("Access-Control-Allow-Origin", origin)
		r.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		r.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		return nil
	}
	return proxy, nil
}

func (h *HttpEngine) EnableProxy(addr string, pattern string) error {
	proxy, err := reverseProxy(addr)
	if err != nil {
		return err
	}
	h.GetChiRouter().Handle(pattern, proxy)
	return nil
}
