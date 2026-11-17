package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	favModels "api-tests-template/internal/managers/favorites/models"
)

// ContainsFavorite проверяет, что объявление есть в списке избранного
func ContainsFavorite(t *testing.T, items []favModels.FavoriteItem, adId string) {
	found := false
	for _, item := range items {
		if item.Advertisement.Id == adId {
			found = true
		}
	}
	require.Truef(t, found, "Объявление %s не найдено в избранном", adId)
}

// NotContainsFavorite проверяет, что объявления нет в списке избранного
func NotContainsFavorite(t *testing.T, items []favModels.FavoriteItem, adId string) {
	for _, item := range items {
		require.NotEqualf(t, adId, item.Advertisement.Id, "Объявление %s всё ещё есть в избранном", adId)
	}
}
