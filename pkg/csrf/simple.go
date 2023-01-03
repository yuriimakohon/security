// Package csrf is simplified version of gorilla/csrf. Only for educational purpose
package csrf

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	gcsrf "github.com/gorilla/csrf"
	"html/template"
	"net/url"
)

const tokenLength = 32

const (
	tokenKey   = "simple.csrf.Token"
	formKey    = "simple.csrf.Form"
	cookieName = "_simple_csrf"
)

var (
	fieldName   = tokenKey
	defaultAge  = 3600 * 12
	headerName  = "X-CSRF-Token"
	safeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}
)

var TemplateTag = "csrfField"

var (
	ErrNoReferer  = errors.New("referer not supplied")
	ErrBadReferer = errors.New("referer invalid")
	ErrNoToken    = errors.New("CSRF token not found in request")
	ErrBadToken   = errors.New("CSRF token invalid")
)

type csrf struct {
	h  http.Handler
	sc *securecookie.SecureCookie
	st store
}

// SimpleProtect provides simple csrf middleware. That's mean DefaultMode for SameSite, always secure mode and "/" path
func SimpleProtect(authKey []byte) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		cs := new(csrf)
		cs.h = h

		if cs.sc == nil {
			cs.sc = securecookie.New(authKey, nil)
			cs.sc.SetSerializer(securecookie.JSONEncoder{})
			cs.sc.MaxAge(defaultAge)
		}

		if cs.st == nil {
			cs.st = &cookieStore{
				name:     cookieName,
				maxAge:   defaultAge,
				secure:   true,
				httpOnly: true,
				sameSite: gcsrf.SameSiteDefaultMode,
				sc:       cs.sc,
			}
		}

		return cs
	}
}

func (cs *csrf) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get or create csrf token from request
	realToken, err := cs.st.Get(r)
	if err != nil || len(realToken) != tokenLength {
		realToken, err = generateRandomBytes(tokenLength)
		if err != nil {
			unauthorizedHandler(w, r, err)
			return
		}

		err = cs.st.Save(realToken, w)
		if err != nil {
			unauthorizedHandler(w, r, err)
			return
		}
	}

	// save token key and csrf form field name in context to retrieve them in TemplateField()
	r = contextSave(r, tokenKey, mask(realToken))
	r = contextSave(r, formKey, fieldName)

	if !contains(safeMethods, r.Method) {
		if r.URL.Scheme == "https" {
			referer, err := url.Parse(r.Referer())
			if err != nil || referer.String() == "" {
				unauthorizedHandler(w, r, ErrNoReferer)
				return
			}

			valid := sameOrigin(r.URL, referer)

			if valid == false {
				unauthorizedHandler(w, r, ErrBadReferer)
				return
			}
		}

		maskedToken, err := cs.requestToken(r)
		if err != nil {
			unauthorizedHandler(w, r, ErrBadToken)
			return
		}

		if maskedToken == nil {
			unauthorizedHandler(w, r, ErrNoToken)
			return
		}

		requestToken := unmask(maskedToken)

		if !compareTokens(requestToken, realToken) {
			unauthorizedHandler(w, r, ErrBadToken)
			return
		}
	}

	w.Header().Add("Vary", "Cookie")

	cs.h.ServeHTTP(w, r)
}

func unauthorizedHandler(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, fmt.Sprintf("%s - %s",
		http.StatusText(http.StatusForbidden), err),
		http.StatusForbidden)
	return
}

func mask(realToken []byte) string {
	otp, err := generateRandomBytes(tokenLength)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(append(otp, xorToken(otp, realToken)...))
}

func unmask(issued []byte) []byte {
	if len(issued) != tokenLength*2 {
		return nil
	}

	otp := issued[tokenLength:]
	masked := issued[:tokenLength]

	return xorToken(otp, masked)
}

func xorToken(a, b []byte) []byte {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}

	res := make([]byte, n)

	for i := 0; i < n; i++ {
		res[i] = a[i] ^ b[i]
	}

	return res
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func contains(vals []string, s string) bool {
	for _, v := range vals {
		if v == s {
			return true
		}
	}

	return false
}

func sameOrigin(a, b *url.URL) bool {
	return a.Scheme == b.Scheme && a.Host == b.Host
}

func (cs *csrf) requestToken(r *http.Request) ([]byte, error) {
	issued := r.Header.Get(headerName)

	if issued == "" {
		issued = r.PostFormValue(fieldName)
	}

	if issued == "" && r.MultipartForm != nil {
		vals := r.MultipartForm.Value[fieldName]

		if len(vals) > 0 {
			issued = vals[0]
		}
	}

	if issued == "" {
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(issued)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func compareTokens(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	return subtle.ConstantTimeCompare(a, b) == 1
}

func TemplateField(r *http.Request) template.HTML {
	if name, err := contextGet(r, formKey); err == nil {
		fragment := fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`,
			name, Token(r))

		return template.HTML(fragment)
	}

	return template.HTML("")
}

func Token(r *http.Request) string {
	if val, err := contextGet(r, tokenKey); err == nil {
		if maskedToken, ok := val.(string); ok {
			return maskedToken
		}
	}

	return ""
}
