package middlewares

import (
	"github.com/Cliengo/acelle-mail/config"
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	conf "github.com/spf13/viper"
	"net/http"
	"strings"
)

var (
	ErrNoTokenFound  = errors.New("no token found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrNotValidToken = errors.New("not valid authorization")
)

type SessionClaims struct {
	jwt.StandardClaims
	CompanyID  string   `json:"company"`
	UserID     string   `json:"user"`
	Privileges []string `json:"privileges"`
}

func JwtAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		stringToken := retrieveToken(r)
		if stringToken == "" {
			http.Error(w, ErrNoTokenFound.Error(), http.StatusUnauthorized)
			return
		}
		token, err := decodeToken(stringToken)
		if err != nil {
			logger.Log.Errorf("%s", err)
			http.Error(w, ErrNotValidToken.Error(), http.StatusForbidden)
			return
		} else if !token.Valid {
			http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*SessionClaims)

		ctx = ContextAppendValues(ctx, claims, SetUserID, SetCompanyID, SetPrivileges)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

var lookupToken = map[string]func(*http.Request) string{
	"header": headerToken,
	"cookie": cookieToken,
	"query":  queryToken,
}

func headerToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func queryToken(r *http.Request) string {
	return r.URL.Query().Get("jwt")
}

func cookieToken(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func retrieveToken(r *http.Request) string {
	for _, fn := range lookupToken {
		if token := fn(r); token != "" {
			return token
		}
	}
	return ""
}

func decodeToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.GetString(config.Secret)), nil
	})
}
