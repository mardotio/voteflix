package utils

import (
	"github.com/go-chi/jwtauth/v5"
	"time"
)

var tokenAuth *jwtauth.JWTAuth

func GetTokenAuth() *jwtauth.JWTAuth {
	if nil == tokenAuth {
		tokenAuth = jwtauth.New("HS256", GetAppConfig().JwtSecret, nil)
		return tokenAuth
	}
	return tokenAuth
}

type JwtClaims struct {
	Sub string
}

func (claims JwtClaims) ToClaimsMap(duration time.Duration) map[string]interface{} {
	claimsMap := map[string]interface{}{
		"sub": claims.Sub,
	}
	jwtauth.SetIssuedNow(claimsMap)
	jwtauth.SetExpiryIn(claimsMap, duration)

	return claimsMap
}
