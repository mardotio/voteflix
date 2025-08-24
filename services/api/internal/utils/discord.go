package utils

import (
	"fmt"
	"voteflix/api/internal/models"
)

func GetAvatarUrl(user models.User) *string {
	if user.DiscordAvatarId == nil {
		return nil
	}

	url := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", user.DiscordId, *user.DiscordAvatarId)

	return &url
}
