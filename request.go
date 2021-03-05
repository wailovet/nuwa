package nuwa

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Request struct {
	GET           map[string]string
	POST          map[string]string
	REQUEST       map[string]string
	COOKIE        map[string]string
	SESSION       map[string]string
	HEADER        map[string]string
	BODY          string
	FILES         map[string][]*multipart.FileHeader
	FILE          *multipart.FileHeader
	OriginRequest *http.Request
	OriginValues  url.Values
}

func (r *Request) SyncGetData(request *http.Request) {
	if r.OriginRequest == nil {
		r.OriginRequest = request
	}
	get := request.URL.Query()

	r.OriginValues = get
	r.GET = make(map[string]string)

	for k := range get {

		str := request.URL.Query().Get(k)
		tmp, err := url.QueryUnescape(str)
		if err != nil {
			log.Println(err.Error())
			r.GET[k] = str
			r.REQUEST[k] = str
		} else {
			r.GET[k] = tmp
			r.REQUEST[k] = tmp
		}
	}
}

func (r *Request) SyncPostData(request *http.Request, mem int64) {
	if r.OriginRequest == nil {
		r.OriginRequest = request
	}
	request.ParseForm()
	request.ParseMultipartForm(mem)
	r.POST = make(map[string]string)

	post := request.PostForm

	for k := range post {
		r.OriginValues[k] = post[k]
		str := request.PostFormValue(k)
		tmp, err := url.QueryUnescape(str)
		if err != nil {
			log.Println(err.Error())
			r.POST[k] = str
			r.REQUEST[k] = str
		} else {
			r.POST[k] = tmp
			r.REQUEST[k] = tmp
		}
	}

	if request.MultipartForm != nil {
		r.FILES = request.MultipartForm.File
		if len(r.FILES) > 0 {
			for fe := range r.FILES {
				for fk := range r.FILES[fe] {
					r.FILE = r.FILES[fe][fk]
				}
			}
		}

		mf := request.MultipartForm.Value
		for k := range mf {
			r.OriginValues[k] = mf[k]
			if len(mf[k]) > 0 {
				r.POST[k] = mf[k][0]
				r.REQUEST[k] = mf[k][0]
			}
		}
	}

	body, _ := ioutil.ReadAll(request.Body)
	r.BODY = string(body)

	rr := map[string]string{}
	json.Unmarshal([]byte(r.BODY), &rr)

	for e := range rr {
		r.POST[e] = rr[e]
		r.REQUEST[e] = rr[e]
	}
}

func (r *Request) SyncHeaderData(request *http.Request) {
	if r.OriginRequest == nil {
		r.OriginRequest = request
	}
	r.HEADER = make(map[string]string)
	header := request.Header
	for k := range header {
		if len(header[k]) > 0 {
			r.HEADER[k] = header[k][0]
		}
	}

}

func (r *Request) SyncCookieData(request *http.Request) {
	if r.OriginRequest == nil {
		r.OriginRequest = request
	}
	cookie := request.Cookies()
	r.COOKIE = make(map[string]string)
	for k := range cookie {
		r.COOKIE[cookie[k].Name] = cookie[k].Value
	}
}

func (r *Request) SyncSessionData(session *session) {
	r.SESSION = session.GetSession()
}
