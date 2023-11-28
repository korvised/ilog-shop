package jwtauth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
	"sync"
)

func NewApiKey(secret string) AuthFactory {
	return &apiKey{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: &Claims{},
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "ilog-shop",
					Subject:   "api-key",
					Audience:  jwt.ClaimStrings{"ilog-shop"},
					ExpiresAt: JwtTimeDurationCal(31560000),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

// ApiKey generator
var apiKeyInstance string
var once sync.Once

func SetApiKey(secret string) {
	once.Do(func() {
		apiKeyInstance = NewApiKey(secret).SingToken()
	})
}

func SetApiKeyInContext(ctx *context.Context) {
	*ctx = metadata.NewOutgoingContext(*ctx, metadata.Pairs("auth", apiKeyInstance))
}
