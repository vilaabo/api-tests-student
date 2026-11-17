package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

// GetFavorites возвращает избранное пользователя, ожидает 200
func GetFavorites(t *testing.T, token string) models.FavoritesList {
	favoritesRawResponse := favorites.HttpGetFavorites(t, token)
	defer favoritesRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(favoritesRawResponse.Body)
	require.NoError(t, errResponse)

	require.Equalf(t, http.StatusOK, favoritesRawResponse.StatusCode,
		"Не удалось получить избранное, ответ сервера: %s", string(respArray))

	favoritesList := models.FavoritesList{}
	err := json.Unmarshal(respArray, &favoritesList)
	require.NoError(t, err)

	return favoritesList
}

// AddFavorite добавляет объявление в избранное, проверяет код ответа и возвращает тело ответа
func AddFavorite(t *testing.T, token, adId string, expectedStatusCode int) string {
	request := models.AddFavoriteRequest{
		AdvertisementId: adId,
	}

	addFavoriteRawResponse := favorites.HttpPostFavorite(t, token, request)
	defer addFavoriteRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(addFavoriteRawResponse.Body)
	require.NoError(t, errResponse)

	require.Equalf(t, expectedStatusCode, addFavoriteRawResponse.StatusCode,
		"Неожиданный код ответа при добавлении в избранное, тело ответа: %s", string(respArray))

	return string(respArray)
}

// RemoveFavorite удаляет объявление из избранного, проверяет код ответа и возвращает тело ответа
func RemoveFavorite(t *testing.T, token, adId string, expectedStatusCode int) string {
	removeFavoriteRawResponse := favorites.HttpDeleteFavorite(t, token, adId)
	defer removeFavoriteRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(removeFavoriteRawResponse.Body)
	require.NoError(t, errResponse)

	require.Equalf(t, expectedStatusCode, removeFavoriteRawResponse.StatusCode,
		"Неожиданный код ответа при удалении из избранного, тело ответа: %s", string(respArray))

	return string(respArray)
}
