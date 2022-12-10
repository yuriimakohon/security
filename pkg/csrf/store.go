package csrf

import (
	"net/http"
	"time"

	gcsrf "github.com/gorilla/csrf"
	"github.com/gorilla/securecookie"
)

type store interface {
	Get(*http.Request) ([]byte, error)
	Save(token []byte, w http.ResponseWriter) error
}

type cookieStore struct {
	name     string
	maxAge   int
	secure   bool
	httpOnly bool
	path     string
	domain   string
	sc       *securecookie.SecureCookie
	sameSite gcsrf.SameSiteMode
}

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
		Path:     cs.path,
		Domain:   cs.domain,
	}

	if cs.maxAge > 0 {
		cookie.Expires = time.Now().Add(
			time.Duration(cs.maxAge) * time.Second)
	}

	http.SetCookie(w, cookie)

	return nil
}
