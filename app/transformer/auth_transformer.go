package transformer

import (
	"go-fiber-starter/app/models"
	"time"
)

type AuthResponse struct {
	UserGUID  string    `json:"user_guid"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func AuthLoginTransformer(user models.User, token string, expiresAt int64) AuthResponse {
	return AuthResponse{
		UserGUID:  user.GUID.String(),
		Name:      user.Name,
		Token:     token,
		ExpiresAt: time.Unix(expiresAt, 0),
	}
}
