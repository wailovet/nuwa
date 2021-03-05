package nuwa

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

type session struct {
	secureCookie *securecookie.SecureCookie
	r            *http.Request
	w            http.ResponseWriter
	data         map[string]string
}

func NewSession(r *http.Request, w http.ResponseWriter) *session {
	session := &session{} // not thread safe
	var hashKey = []byte("osmanthuswine-very-secret")
	// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
	// Shorter keys may weaken the encryption used.
	var blockKey = []byte("osmanthuswine-lot-secret")
	session.secureCookie = securecookie.New(hashKey, blockKey)
	session.r = r
	session.w = w

	session.data = make(map[string]string)
	if cookie, err := session.r.Cookie("osmseccidhas"); err == nil {
		session.secureCookie.Decode("osmseccidhas", cookie.Value, &session.data)
	} else {
		//helper.GetInstanceLog().Out(err.Error())
	}
	return session
}

func (session *session) GetSession() map[string]string {
	return session.data
}

func (session *session) SetSession(value map[string]string) {
	if encoded, err := session.secureCookie.Encode("osmseccidhas", value); err == nil {
		cookie := &http.Cookie{
			Name:     "osmseccidhas",
			Value:    encoded,
			Path:     "/",
			Secure:   false,
			HttpOnly: false,
		}
		http.SetCookie(session.w, cookie)
		session.data = value
	}
}
