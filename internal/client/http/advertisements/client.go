package advertisements

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpGetAdvertisements выполняет GET-запрос поиска объявлений по тексту
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Query("limit", "50"). // максимум по swagger, страхует от пагинации
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpPostAdvertisement выполняет POST-запрос создания объявления с полями формы и фотографиями
func HttpPostAdvertisement(t *testing.T, token string, fields map[string]string, photos []string) *http.Response {
	request := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath)
	for name, value := range fields {
		request = request.MultipartFormData(name, value)
	}
	if len(photos) > 0 {
		request = request.MultipartFile(models.FieldPhotos, photos...)
	}
	return request.Expect(t).End().Response
}

// HttpGetAdvertisementById выполняет GET-запрос получения объявления по id
func HttpGetAdvertisementById(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement выполняет DELETE-запрос удаления объявления по id
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
