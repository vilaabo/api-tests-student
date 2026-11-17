package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpGetFavorites выполняет GET-запрос получения избранного пользователя
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page_size", "50"). // максимум по swagger, в избранном тоже не больше 50 объявлений
		Expect(t).
		End().Response
}

// HttpPostFavorite выполняет POST-запрос добавления объявления в избранное
func HttpPostFavorite(t *testing.T, token string, request models.AddFavoriteRequest) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpPostFavoritesBulk выполняет POST-запрос добавления в избранное сразу нескольких объявлений
func HttpPostFavoritesBulk(t *testing.T, token string, request models.BulkAddFavoritesRequest) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesBulkPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite выполняет DELETE-запрос удаления объявления из избранного
func HttpDeleteFavorite(t *testing.T, token, adId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.FavoritesPath + "/" + adId).
		Expect(t).
		End().Response
}
