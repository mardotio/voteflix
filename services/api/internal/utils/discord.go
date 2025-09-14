package utils

import (
	"fmt"
)

func GetAvatarUrl(discordUserId string, avatarId *string, isServer bool) *string {
	if avatarId == nil {
		return nil
	}

	path := "avatars"

	if isServer {
		path = "icons"
	}

	url := fmt.Sprintf("https://cdn.discordapp.com/%s/%s/%s.png", path, discordUserId, *avatarId)

	return &url
}
