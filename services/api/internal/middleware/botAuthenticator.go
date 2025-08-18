package middleware

import (
	"fmt"
	"net/http"
	"voteflix/api/internal/app"
	"voteflix/api/internal/utils"
)

func BotAuthenticator(config *app.AppConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			jsonSender := utils.NewJsonSender(w, r)
			discordBotKey := r.Header.Get("Bot-Api-Key")

			if discordBotKey == "" {
				jsonSender.Unauthorized(fmt.Errorf("no bot token found"))
				return
			}

			if discordBotKey != config.ApiBotKey {
				jsonSender.Unauthorized(fmt.Errorf("bot not authorized"))
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(handler)
	}
}
