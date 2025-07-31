package auth

import (
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strconv"
	"time"
	"voteflix/api/src/utils"
)

type jsonEpochTime time.Time

func (t jsonEpochTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

type createTokenResponse struct {
	Token   string        `json:"token"`
	Expires jsonEpochTime `json:"expires"`
}

func (rd createTokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	claims := utils.JwtClaims{Sub: "user-id"}
	token, tokenString, _ := utils.GetTokenAuth().Encode(claims.ToClaimsMap(time.Hour * 1))

	render.Status(r, http.StatusOK)
	err := render.Render(w, r, createTokenResponse{Token: tokenString, Expires: jsonEpochTime(token.Expiration())})

	if nil != err {
		log.Println(err)
	}
}
