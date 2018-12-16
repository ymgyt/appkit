package middlewares

import (
	"net/http"

	"github.com/ymgyt/appkit/services"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

// JWTVerifier -
type JWTVerifier struct {
	logger     *zap.Logger
	hmacSecret []byte
}

// JWTVerifyConfig -
type JWTVerifyConfig struct {
	Logger     *zap.Logger
	HMACSecret []byte
}

// MustJWTVerifier -
func MustJWTVerifier(cfg *JWTVerifyConfig) *JWTVerifier {
	j, err := NewJWTVerifier(cfg)
	if err != nil {
		panic(err)
	}
	return j
}

// NewJWTVerifier -
func NewJWTVerifier(cfg *JWTVerifyConfig) (*JWTVerifier, error) {
	return &JWTVerifier{logger: cfg.Logger, hmacSecret: cfg.HMACSecret}, nil
}

// ServeHTTP -
func (m *JWTVerifier) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rawToken := m.readToken(r)
	if rawToken == "" {
		m.logger.Debug("jwt token nof found")
		next(w, r)
		return
	}

	token, err := jwt.Parse(rawToken, m.keyFunc)
	if err != nil {
		// m.next.ServeHTTP(w, r)
		m.logger.Debug("jwt parse failed", zap.String("err", err.Error()))
		next(w, r)
		return
	}
	// m.logger.Debug("jwt parse success", zap.Any("claims", token.Claims))

	rr := r.WithContext(services.SetIDToken(r.Context(), token))
	// m.next.ServeHTTP(w, rr)
	next(w, rr)
}

func (m *JWTVerifier) keyFunc(token *jwt.Token) (interface{}, error) {
	// TODO more verify
	return m.hmacSecret, nil
}

func (m *JWTVerifier) readToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token != "" {
		return token
	}

	return r.URL.Query().Get("id_token")
}
