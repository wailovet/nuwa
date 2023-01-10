package nuwa

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
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
	displayCallback      func(data []byte, code int)
	code                 int
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
	if r.displayCallback != nil {
		r.displayCallback(data, r.code)
	}
	panic(nil)
}

func (r *Response) DisplayByRawCache(data []byte, code int) {

	cc := Config()
	//log.Println("crossDomain:", cc.CrossDomain)
	if cc.CrossDomain != "" {
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Origin", cc.CrossDomain)
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Methods", "Access-Control-Allow-Methods")
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With")
	}
	if code > 0 {
		r.OriginResponseWriter.WriteHeader(code)
	}
	_, _ = r.OriginResponseWriter.Write(data)
	panic(nil)
}

func (r *Response) DisplayCallback(displayCallback func(data []byte, code int)) {
	r.displayCallback = displayCallback
}

func (r *Response) DisplayByString(data string) {
	r.DisplayByRaw([]byte(data))
}

func (r *Response) Display(data interface{}, msg string, code int) {
	result := responseData{code, data, msg}
	text, err := json.Marshal(result)
	if err != nil {
		r.OriginResponseWriter.WriteHeader(500)
		r.code = 500
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
		r.DisplayByError(err.Error(), 504)
	}
}

func (r *Response) DisplayBySuccess(msgs ...string) {
	msg := "success"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
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

func (r *Response) DisplayJPEG(img image.Image, o ...*jpeg.Options) {

	cc := Config()
	if cc.CrossDomain != "" {
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Origin", cc.CrossDomain)
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Methods", "Access-Control-Allow-Methods")
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With")
	}

	r.OriginResponseWriter.Header().Set("Content-Type", "image/jpeg")
	opt := &jpeg.Options{}
	if len(o) > 0 {
		opt = o[0]
	}

	jpeg.Encode(r.OriginResponseWriter, img, opt)
	panic(nil)
}
func (r *Response) DisplayPNG(img image.Image) {

	cc := Config()
	if cc.CrossDomain != "" {
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Origin", cc.CrossDomain)
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Methods", "Access-Control-Allow-Methods")
		r.OriginResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With")
	}

	r.OriginResponseWriter.Header().Set("Content-Type", "image/png")

	png.Encode(r.OriginResponseWriter, img)
	panic(nil)
}
