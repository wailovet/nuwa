package nuwa

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strings"
)

type responseData struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type Response struct {
	Session              *session
	OriginResponseWriter http.ResponseWriter
}

func (r *Response) DisplayByRaw(data []byte) {

	cc := Config()
	//log.Println("crossDomain:", cc.CrossDomain)
	if cc.CrossDomain != "" {
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Origin", cc.CrossDomain)
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Methods", "Access-Control-Allow-Methods")
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With")
	}
	_, _ = r.OriginResponseWriter.Write(data)
	panic(nil)
}

func (r *Response) DisplayByString(data string) {
	r.DisplayByRaw([]byte(data))
}

func (r *Response) Display(data interface{}, msg string, code int) {
	result := responseData{code, data, msg}
	text, err := json.Marshal(result)
	if err != nil {
		r.OriginResponseWriter.WriteHeader(500)
		r.DisplayByString("服务器异常:" + err.Error())
	}
	r.DisplayByRaw(text)
}

func (r *Response) DisplayByError(msg string, code int, data ...string) {
	result := responseData{code, data, msg}
	text, err := json.Marshal(result)
	if err != nil {
		r.Display(nil, "JSON返回格式解析异常:"+err.Error(), 500)
	}
	r.DisplayByRaw(text)
}

func (r *Response) CheckErrDisplayByError(err error, msg ...string) {
	if err == nil {
		return
	}
	if len(msg) > 0 {
		r.DisplayByError(strings.Join(msg, ","), 504)
	} else {
		r.DisplayByError(err.Error(), 504, strings.Split(string(debug.Stack()), "\n\t")...)
	}
}

func (r *Response) DisplayBySuccess(msg string) {
	result := responseData{0, nil, msg}
	text, err := json.Marshal(result)
	if err != nil {
		r.Display(nil, "JSON返回格式解析异常:"+err.Error(), 500)
	}
	r.DisplayByRaw(text)
}

func (r *Response) DisplayByData(data interface{}) {
	result := responseData{0, data, ""}
	text, err := json.Marshal(result)
	if err != nil {
		r.Display(nil, "JSON返回格式解析异常:"+err.Error(), 500)
	}
	r.DisplayByRaw(text)
}

func (r *Response) SetSession(name string, value string) {
	data := r.Session.GetSession()
	data[name] = value
	r.Session.SetSession(data)
}

func (r *Response) DeleteSession(name string) {
	data := r.Session.GetSession()
	delete(data, name)
	r.Session.SetSession(data)
}

func (r *Response) ClearSession() {
	data := make(map[string]string)
	r.Session.SetSession(data)
}

func (r *Response) SetCookie(name string, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
	}
	http.SetCookie(r.OriginResponseWriter, cookie)
}

func (r *Response) SetHeader(name string, value string) {
	r.OriginResponseWriter.Header().Set(name, value)
}
