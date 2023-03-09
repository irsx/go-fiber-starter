package transformer

import (
	"go-fiber-starter/app/models"
	"time"
)

type NewsDetailResponse struct {
	GUID        string    `json:"guid"`
	UserGUID    string    `json:"user_guid"`
	Title       string    `json:"title"`
	Image       string    `json:"image"`
	Description string    `json:"description,omitempty"`
	Status      bool      `json:"status"`
	HyperLink   string    `json:"hyperlink"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewsDetailTransformer(news models.News) (detail *NewsDetailResponse) {
	return &NewsDetailResponse{
		GUID:        news.GUID.String(),
		UserGUID:    news.UserGUID,
		Title:       news.Title,
		Image:       news.Image,
		Description: news.Description,
		HyperLink:   news.HyperLink,
		Status:      news.Status == "1",
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
	}
}

func NewsListTransformer(newsData []*models.News) []NewsDetailResponse {
	list := make([]NewsDetailResponse, 0)
	for _, news := range newsData {
		list = append(list, NewsDetailResponse{
			GUID:      news.GUID.String(),
			UserGUID:  news.UserGUID,
			Title:     news.Title,
			Image:     news.Image,
			HyperLink: news.HyperLink,
			Status:    news.Status == "1",
			CreatedAt: news.CreatedAt,
			UpdatedAt: news.UpdatedAt,
		})
	}

	return list
}
