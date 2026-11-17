package models

import (
	advModels "api-tests-template/internal/managers/advertisements/models"
)

type FavoritesList struct {
	Items      []FavoriteItem `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}

type BulkAddFavoritesResult struct {
	Added        []string `json:"added"`
	Skipped      []string `json:"skipped"`
	NotFound     []string `json:"not_found"`
	LimitReached bool     `json:"limit_reached"`
}

type FavoriteItem struct {
	FavoriteId    string                  `json:"favorite_id"`
	AddedAt       string                  `json:"added_at"`
	Advertisement advModels.Advertisement `json:"advertisement"`
}
