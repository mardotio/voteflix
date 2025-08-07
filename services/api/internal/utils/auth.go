package utils

import (
	"github.com/go-chi/jwtauth/v5"
	"time"
)

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
