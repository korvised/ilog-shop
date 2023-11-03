package jwtauth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/korvised/ilog-shop/pkg/utils"
	"math"
	"time"
)

type (
	AuthFactory interface {
		SingToken() string
	}

	Claims struct {
		ID       string `json:"id"`
		RoleCode int    `json:"role_code"`
	}

	AuthMapClaims struct {
		*Claims
		jwt.RegisteredClaims
	}

	authConcrete struct {
		Secret []byte
		Claims *AuthMapClaims `json:"claims"`
	}

	accessToken  struct{ *authConcrete }
	refreshToken struct{ *authConcrete }
	apiKey       struct{ *authConcrete }
)

func NewAuthFactory(secret []byte, claims *AuthMapClaims) AuthFactory {
	return &accessToken{&authConcrete{secret, claims}}
}

func (a *authConcrete) SingToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.Claims)
	ss, _ := token.SignedString(a.Secret)
	return ss
}

// JwtTimeDurationCal t is a second unit
func JwtTimeDurationCal(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(now().Add(time.Duration(t * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func now() time.Time { return time.Now().In(utils.LoadLocation()) }

func ReloadToken(secret string, expiredAt int64, claims *Claims) string {
	obj := &refreshToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "ilog-shop",
					Subject:   "reload-token",
					Audience:  jwt.ClaimStrings{"ilog-shop"},
					ExpiresAt: jwtTimeRepeatAdapter(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}

	return obj.SingToken()
}

func ParseToken(secret string, tokenString string) (*AuthMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error: unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.New("error: token format is invalid")
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errors.New("error: token is expired")
		default:
			return nil, errors.New("error: token is invalid")
		}
	}

	if claim, ok := token.Claims.(*AuthMapClaims); ok {
		return claim, nil
	} else {
		return nil, errors.New("error: claim type is invalid")
	}
}
