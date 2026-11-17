package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/messages"
	advManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	favManager "api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

const (
	favoritesLimit = 50
)

// Negative: повторное добавление того же объявления в избранное
func (s *TestSuite) TestAddFavoriteTwice() {
	advertisement := s.createSellerAdvertisement()
	defer advManager.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)

	base.Step("Добавляем объявление в избранное первый раз")
	favManager.AddFavorite(s.T(), s.buyer.Token, advertisement.Id, http.StatusCreated)

	base.Step("Добавляем то же объявление второй раз")
	body := favManager.AddFavorite(s.T(), s.buyer.Token, advertisement.Id, http.StatusConflict)

	assertions.CheckError(s.T(), body, messages.ErrorAlreadyExists, messages.AlreadyInFavorites)
	base.Check("Получили 409 и сообщение «" + messages.AlreadyInFavorites + "»")
}

// Negative: добавление в избранное своего объявления
func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	advertisement := s.createSellerAdvertisement()
	defer advManager.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)

	base.Step("Продавец добавляет в избранное своё объявление")
	body := favManager.AddFavorite(s.T(), s.seller.Token, advertisement.Id, http.StatusForbidden)

	assertions.CheckError(s.T(), body, messages.ErrorForbidden, messages.CannotAddOwnAdvertisement)
	base.Check("Получили 403 и сообщение «" + messages.CannotAddOwnAdvertisement + "»")
}

// Negative: добавление в избранное без advertisement_id
func (s *TestSuite) TestAddFavoriteWithoutAdvertisementId() {
	base.Step("Отправляем запрос на добавление в избранное без advertisement_id")
	body := favManager.AddFavorite(s.T(), s.buyer.Token, "", http.StatusBadRequest)

	assertions.CheckError(s.T(), body, messages.ErrorInvalidRequest, messages.AdvertisementIdRequired)
	base.Check("Получили 400 и сообщение, что advertisement_id обязателен")
}

// Negative: в избранное нельзя добавить больше 50 объявлений
func (s *TestSuite) TestFavoritesLimit() {
	base.Precondition("Регистрируем отдельного покупателя, чтобы не заполнять избранное основного аккаунта")
	newBuyer := auth.RegisterRandomUser(s.T(), "Buyer")

	base.Precondition("Берём 50 самых старых объявлений из ленты, они не пропадут во время теста")
	oldestAdvertisements := advManager.GetOldestAdvertisements(s.T(), newBuyer.Token)
	require.Len(s.T(), oldestAdvertisements.Items, favoritesLimit, "В ленте меньше 50 объявлений")

	advertisement := s.createSellerAdvertisement()
	defer advManager.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)

	s.Run("Добавляем 50 объявлений", func() {
		adIds := []string{}
		for _, item := range oldestAdvertisements.Items {
			adIds = append(adIds, item.Id)
		}

		// одним запросом, потому что на 50 запросов подряд сервер отвечает 503
		base.Step("Добавляем в избранное 50 объявлений через /favorites/bulk")
		bulkResult := favManager.AddFavoritesBulk(s.T(), newBuyer.Token, adIds)
		require.Len(s.T(), bulkResult.Added, favoritesLimit, "Добавились не все 50 объявлений")

		favorites := favManager.GetFavorites(s.T(), newBuyer.Token)
		require.Equal(s.T(), favoritesLimit, favorites.Total, "В избранном должно быть 50 объявлений")
		base.Check("В избранном 50 объявлений")
	})

	s.Run("51-е объявление не добавляется", func() {
		base.Step("Добавляем в избранное 51-е объявление")
		body := favManager.AddFavorite(s.T(), newBuyer.Token, advertisement.Id, http.StatusBadRequest)

		assertions.CheckError(s.T(), body, messages.ErrorLimitReached, messages.FavoritesLimitReached)
		base.Check("Получили 400 и сообщение «" + messages.FavoritesLimitReached + "»")
	})
}
