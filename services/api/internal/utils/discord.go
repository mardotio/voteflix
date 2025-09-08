package utils

import (
	"fmt"
)

func GetAvatarUrl(discordUserId string, avatarId *string) *string {
	if avatarId == nil {
		return nil
	}

	url := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", discordUserId, *avatarId)

	return &url
}
