package jwt

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWTAuth struct {
	algorithm       jwa.SignatureAlgorithm
	signKey         interface{}
	verifyKey       interface{}
	verifier        jwt.ParseOption
	validateOptions []jwt.ValidateOption
}

var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}

	ErrInvalidToken     = errors.New("invalid token")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrExpired          = errors.New("token expired")
	ErrNBFInvalid       = errors.New("token not valid yet")
	ErrIATInvalid       = errors.New("token issued at time is invalid")
	ErrNoTokenFound     = errors.New("no token found")
	ErrAlgorithmInvalid = errors.New("token algorithm is invalid")
)

func New(alg string, signKey interface{}, verifyKey interface{}, validateOptions ...jwt.ValidateOption) *JWTAuth {
	ja := &JWTAuth{
		algorithm:       jwa.SignatureAlgorithm(alg),
		signKey:         signKey,
		verifyKey:       verifyKey,
		validateOptions: validateOptions,
	}

	if ja.verifyKey != nil {
		ja.verifier = jwt.WithKey(ja.algorithm, ja.verifyKey)
	} else {
		ja.verifier = jwt.WithKey(ja.algorithm, ja.signKey)
	}

	return ja
}

// Verifier middleware will verify a JWT string from an HTTP request.
// Will look for a token in the Authorization header, must be passed as a
// Bearer token.
func Verifier(ja *JWTAuth) func(next http.Handler) http.Handler {
	return Verify(ja, TokenFromHeader)
}

func Verify(ja *JWTAuth, findTokenFns ...func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfh := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := VerifyRequest(ja, r, findTokenFns...)
			ctx = NewContext(ctx, token, err)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(hfh)
	}
}

func VerifyRequest(ja *JWTAuth, r *http.Request, findTokenFns ...func(r *http.Request) string) (jwt.Token, error) {
	var tokenString string

	for _, fn := range findTokenFns {
		tokenString = fn(r)
		if tokenString != "" {
			break
		}
	}
	if tokenString == "" {
		return nil, ErrNoTokenFound
	}

	return VerifyToken(ja, tokenString)
}

func VerifyToken(ja *JWTAuth, tokenString string) (jwt.Token, error) {
	token, err := ja.Decode(tokenString)
	if err != nil {
		return token, ErrorReason(err)
	}

	if token == nil {
		return nil, ErrUnauthorized
	}

	if err := jwt.Validate(token, ja.validateOptions...); err != nil {
		return token, ErrorReason(err)
	}

	return token, nil
}

func (ja *JWTAuth) Encode(claims map[string]interface{}) (t jwt.Token, tokenString string, err error) {
	t = jwt.New()
	for k, v := range claims {
		t.Set(k, v)
	}
	payload, err := ja.sign(t)
	if err != nil {
		return nil, "", err
	}
	tokenString = string(payload)
	return
}

func (ja *JWTAuth) Decode(tokenString string) (jwt.Token, error) {
	return ja.parse([]byte(tokenString))
}

func (ja *JWTAuth) ValidateOptions() []jwt.ValidateOption {
	return ja.validateOptions
}

func (ja *JWTAuth) sign(token jwt.Token) ([]byte, error) {
	return jwt.Sign(token, jwt.WithKey(ja.algorithm, ja.signKey))
}

func (ja *JWTAuth) parse(payload []byte) (jwt.Token, error) {
	return jwt.Parse(payload, ja.verifier, jwt.WithValidate(false))
}

func ErrorReason(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired()), err == ErrExpired:
		return ErrExpired
	case errors.Is(err, jwt.ErrInvalidIssuedAt()), err == ErrIATInvalid:
		return ErrIATInvalid
	case errors.Is(err, jwt.ErrTokenNotYetValid()), err == ErrNBFInvalid:
		return ErrNBFInvalid
	default:
		return ErrUnauthorized
	}
}

func Authenticator(ja *JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := FromContext(r.Context())

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if token == nil || jwt.Validate(token, ja.validateOptions...) != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

func NewContext(ctx context.Context, t jwt.Token, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, t)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}

func FromContext(ctx context.Context) (jwt.Token, map[string]interface{}, error) {
	token, _ := ctx.Value(TokenCtxKey).(jwt.Token)

	var err error
	var claims map[string]interface{}

	if token != nil {
		claims, err = token.AsMap(context.Background())
		if err != nil {
			return token, nil, err
		}
	} else {
		claims = map[string]interface{}{}
	}

	err, _ = ctx.Value(ErrorCtxKey).(error)

	return token, claims, err
}

func UnixTime(tm time.Time) int64 {
	return tm.UTC().Unix()
}

func EpochNow() int64 {
	return time.Now().UTC().Unix()
}

func ExpireIn(tm time.Duration) int64 {
	return EpochNow() + int64(tm.Seconds())
}

func SetIssuedAt(claims map[string]interface{}, tm time.Time) {
	claims["iat"] = tm.UTC().Unix()
}

func SetIssuedNow(claims map[string]interface{}) {
	claims["iat"] = EpochNow()
}

func SetExpiry(claims map[string]interface{}, tm time.Time) {
	claims["exp"] = tm.UTC().Unix()
}

func SetExpiryIn(claims map[string]interface{}, tm time.Duration) {
	claims["exp"] = ExpireIn(tm)
}

func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func TokenFromQuery(r *http.Request) string {
	// Get token from query param named "jwt".
	return r.URL.Query().Get("jwt")
}

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}
