package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

// ContainsAdvertisementByTitle проверяет, что в списке есть объявление с указанным названием
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	titles := make([]string, 0, len(items))
	for _, item := range items {
		if item.Title == title {
			return
		}
		titles = append(titles, item.Title)
	}
	require.Failf(t, "Объявление не найдено",
		"Ожидали объявление %q, а в ответе: %q", title, titles)
}
