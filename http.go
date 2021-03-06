package nuwa

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewHttp(config *config) *HttpEngine {
	return &HttpEngine{
		InstanceConfig: config,
	}
}

func Http() *HttpEngine {
	http_.InstanceConfig = Config()
	return &http_
}

var http_ = HttpEngine{}

type HttpEngine struct {
	chiRouter      *chi.Mux
	isDebug        bool
	InstanceConfig *config
}

func (h *HttpEngine) DisableDebug() {
	h.isDebug = false
}

func (h *HttpEngine) GetChiRouter() *chi.Mux {
	if h.chiRouter == nil {
		h.chiRouter = chi.NewRouter()
		if h.isDebug {
			h.chiRouter.Use(middleware.Logger)
		}

		h.chiRouter.Use(middleware.RequestID)
		h.chiRouter.Use(middleware.RealIP)
		h.chiRouter.Use(middleware.Recoverer)
	}
	return h.chiRouter
}

func (h *HttpEngine) Run() error {
	cc := h.InstanceConfig
	Loger().Out("开始监听:", cc.Host+":"+cc.Port)

	r := h.GetChiRouter()

	if cc.StaticRouter == "" {
		cc.StaticRouter = "/*"
	}
	_, err := os.Stat("./html")

	if cc.StaticFileSystem == nil {
		cc.StaticFileSystem = http.Dir("static")
	}

	staticHandle := http.FileServer(cc.StaticFileSystem)
	if cc.IsStaticStripPrefix {
		staticHandle = http.StripPrefix(strings.Replace(cc.StaticRouter, "*", "", -1), staticHandle)
	}
	r.Handle(cc.StaticRouter, staticHandle)

	listener, err := net.Listen("tcp", cc.Host+":"+cc.Port)
	if err != nil {
		return err
	}
	return http.Serve(listener, r)
}

func (h *HttpEngine) HandleFunc(pattern string, callback func(ctx HttpContext)) {

	cc := h.InstanceConfig
	r := h.GetChiRouter()
	r.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {

		requestData := Request{}

		sessionMan := NewSession(request, writer)

		requestData.REQUEST = make(map[string]string)
		//GET
		requestData.SyncGetData(request)
		//POST
		requestData.SyncPostData(request, cc.PostMaxMemory)
		//HEADER
		requestData.SyncHeaderData(request)
		//COOKIE
		requestData.SyncCookieData(request)
		//SESSION
		requestData.SyncSessionData(sessionMan)

		responseHandle := Response{OriginResponseWriter: writer, Session: sessionMan}

		responseHandle.OriginResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
		if cc.CrossDomain != "" {
			responseHandle.OriginResponseWriter.Header().Set("Access-Control-Allow-Origin", cc.CrossDomain)
		}

		defer func() {
			// println("h.isDebug:", h.isDebug)
			errs := recover()
			if errs == nil {
				return
			}
			errtxt := fmt.Sprintf("%v", errs)
			if errtxt != "" {
				if h.isDebug {
					responseHandle.DisplayByError(errtxt, 500, strings.Split(string(debug.Stack()), "\n\t")...)
				} else {
					responseHandle.DisplayByError(errtxt, 500, "service error")
				}
			}
		}()

		callback(HttpContext{
			Request:  requestData,
			Response: responseHandle,
		})
	})
}

func (h *HttpEngine) EnableAuthenticate(user, password string) {
	h.GetChiRouter().Use(func(handler http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !checkAuth(r, user, password) {
				w.Header().Set("WWW-Authenticate", `Basic`)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			} else {
				handler.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	})
}

func checkAuth(r *http.Request, user, pass string) bool {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	return pair[0] == user && pair[1] == pass
}
