package csrf

import (
	"net/http"
	"time"

	gcsrf "github.com/gorilla/csrf"
	"github.com/gorilla/securecookie"
)

// store interface manage storage of csrf-token
type store interface {
	Get(*http.Request) ([]byte, error)
	Save(token []byte, w http.ResponseWriter) error
}

// cookieStore is default(and only one, because store is not exported) cookie-based implementation of store
type cookieStore struct {
	name     string
	maxAge   int
	secure   bool
	httpOnly bool
	sc       *securecookie.SecureCookie
	sameSite gcsrf.SameSiteMode
}

// Get decoded csrf-token from SecureCookie
func (cs *cookieStore) Get(r *http.Request) ([]byte, error) {
	cookie, err := r.Cookie(cs.name)
	if err != nil {
		return nil, err
	}

	token := make([]byte, tokenLength)
	err = cs.sc.Decode(cs.name, cookie.Value, &token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Save encoded token to SecureCookie
func (cs *cookieStore) Save(token []byte, w http.ResponseWriter) error {
	encoded, err := cs.sc.Encode(cs.name, token)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     cs.name,
		Value:    encoded,
		MaxAge:   cs.maxAge,
		HttpOnly: cs.httpOnly,
		Secure:   cs.secure,
		SameSite: http.SameSite(cs.sameSite),
	}

	if cs.maxAge > 0 {
		cookie.Expires = time.Now().Add(
			time.Duration(cs.maxAge) * time.Second)
	}

	http.SetCookie(w, cookie)

	return nil
}
