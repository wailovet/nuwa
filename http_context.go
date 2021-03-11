package nuwa

import (
	"fmt"
	"strconv"
)

type HttpContext struct {
	Request
	Response
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
