package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

// GetAdvertisementsList возвращает список объявлений по поисковому запросу, ожидает 200
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := advertisements.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// CreateAdvertisement создаёт объявление и возвращает его, ожидает 201
func CreateAdvertisement(t *testing.T, token string, fields map[string]string, photos []string) models.Advertisement {
	createAdvertisementRawResponse := advertisements.HttpPostAdvertisement(t, token, fields, photos)
	defer createAdvertisementRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(createAdvertisementRawResponse.Body)
	require.NoError(t, errResponse)

	require.Equalf(t, http.StatusCreated, createAdvertisementRawResponse.StatusCode,
		"Объявление не создалось, ответ сервера: %s", string(respArray))

	advertisement := models.Advertisement{}
	err := json.Unmarshal(respArray, &advertisement)
	require.NoError(t, err)

	return advertisement
}

// GetAdvertisementById возвращает объявление по id, ожидает 200
func GetAdvertisementById(t *testing.T, id string) models.Advertisement {
	advertisementRawResponse := advertisements.HttpGetAdvertisementById(t, id)
	defer advertisementRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(advertisementRawResponse.Body)
	require.NoError(t, errResponse)

	require.Equalf(t, http.StatusOK, advertisementRawResponse.StatusCode,
		"Не удалось получить объявление %s, ответ сервера: %s", id, string(respArray))

	advertisement := models.Advertisement{}
	err := json.Unmarshal(respArray, &advertisement)
	require.NoError(t, err)

	return advertisement
}

// GetAdvertisementPhotos возвращает фотографии объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	photosRawResponse := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	defer photosRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(photosRawResponse.Body)
	require.NoError(t, errResponse)

	require.Equalf(t, http.StatusOK, photosRawResponse.StatusCode,
		"Не удалось получить фото объявления %s, ответ сервера: %s", id, string(respArray))

	photos := []models.Photo{}
	err := json.Unmarshal(respArray, &photos)
	require.NoError(t, err)

	return photos
}

// DeleteAdvertisement удаляет объявление по id, ожидает 204
func DeleteAdvertisement(t *testing.T, token, id string) {
	deleteAdvertisementRawResponse := advertisements.HttpDeleteAdvertisement(t, token, id)
	defer deleteAdvertisementRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(deleteAdvertisementRawResponse.Body)
	require.NoError(t, errResponse)

	require.Equalf(t, http.StatusNoContent, deleteAdvertisementRawResponse.StatusCode,
		"Не удалось удалить объявление %s, ответ сервера: %s", id, string(respArray))
}
