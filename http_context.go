package nuwa

import (
	"fmt"
	"strconv"

	"github.com/xujiajun/nutsdb/inmemory"
)

type HttpContext struct {
	Request
	Response
}

var httpCache = NewNutsdbMemory(inmemory.DefaultOptions)

func (r *HttpContext) EnableCache(ttl uint32) {
	key := Helper().Md5(Helper().JsonEncode([]interface{}{
		r.OriginRequest.URL.Path, r.REQUEST,
	}))
	r.DisplayCallback(func(data []byte, code int) {
		httpCache.Set(key, map[string]interface{}{
			"data": string(data),
			"code": code,
		}, ttl)
	})
	ret := httpCache.Get(key)
	if ret.Map()["data"].String() != "" {
		r.DisplayByRawCache([]byte(ret.Map()["data"].String()), int(ret.Map()["code"].Int()))
	}
}

func (r *HttpContext) ParamRequired(key string) string {
	if r.REQUEST[key] == "" {
		r.DisplayByError(fmt.Sprintf("参数错误,[%s]不允许为空", key), 404)
	}
	return r.REQUEST[key]
}

func (r *HttpContext) ParamRequired2Int(key string) int {
	s := r.ParamRequired(key)
	i, err := strconv.ParseInt(s, 10, 32)
	r.CheckErrDisplayByError(err)
	return int(i)
}

func (r *HttpContext) ParamRequired2Int64(key string) int64 {
	s := r.ParamRequired(key)
	i, err := strconv.ParseInt(s, 10, 64)
	r.CheckErrDisplayByError(err)
	return i
}

func (r *HttpContext) ParamRequired2Float(key string) float64 {
	s := r.ParamRequired(key)
	i, err := strconv.ParseFloat(s, 64)
	r.CheckErrDisplayByError(err)
	return i
}
