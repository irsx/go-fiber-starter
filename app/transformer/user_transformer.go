package transformer

import (
	"encoding/json"
	"go-fiber-starter/app/models"
)

type UserResponse struct {
	GUID        string `json:"guid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func UserListTransformer(user []*models.User) (rows []UserResponse) {
	for _, row := range user {
		var mapResponse UserResponse
		jsonResponse, _ := json.Marshal(row)
		json.Unmarshal(jsonResponse, &mapResponse)

		mapResponse = UserResponse{
			GUID:        row.GUID.String(),
			Name:        row.Name,
			Email:       row.Email,
			PhoneNumber: row.PhoneNumber,
			CreatedAt:   row.CreatedAt.String(),
			UpdatedAt:   row.UpdatedAt.String(),
		}
		rows = append(rows, mapResponse)
	}
	return rows
}
