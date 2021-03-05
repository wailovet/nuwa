package nuwa

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

func NewHttp(config *config) *httpImp {
	http_.InstanceConfig = config
	return &http_
}

func Http() *httpImp {
	http_.InstanceConfig = Config()
	return &http_
}

var http_ = httpImp{}

type httpImp struct {
	chiRouter      *chi.Mux
	isDebug        bool
	InstanceConfig *config
}

func (h *httpImp) Debug() {
	h.isDebug = false
}

func (h *httpImp) GetChiRouter() *chi.Mux {
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

func (h *httpImp) Run() error {
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

func (h *httpImp) HandleFunc(pattern string, callback func(request Request, response Response)) {

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
			errs := recover()
			if errs == nil {
				return
			}
			errtxt := fmt.Sprintf("%v", errs)
			if errtxt != "" {
				responseHandle.DisplayByError(errtxt, 500, strings.Split(string(debug.Stack()), "\n\t")...)
			}
		}()

		callback(requestData, responseHandle)
	})
}
