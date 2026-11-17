package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/messages"
	advManager "api-tests-template/internal/managers/advertisements"
	favManager "api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

const (
	nonExistentAdId    = "7b0a4c56-1d2e-4f3a-9b8c-5d6e7f8a9b0c"
	invalidAdId        = "abc"
	invalidUuidMessage = "invalid UUID length: 3"
)

// Negative: добавление в избранное и удаление из избранного несуществующего объявления
func (s *TestSuite) TestFavoritesWithNonExistentAdvertisement() {
	s.Run("Добавление несуществующего объявления", func() {
		base.Step("Добавляем в избранное объявление с несуществующим id " + nonExistentAdId)
		body := favManager.AddFavorite(s.T(), s.buyer.Token, nonExistentAdId, http.StatusNotFound)

		assertions.CheckError(s.T(), body, messages.ErrorNotFound, messages.AdvertisementNotFound)
		base.Check("Получили 404 и сообщение «" + messages.AdvertisementNotFound + "»")
	})

	s.Run("Удаление несуществующего объявления", func() {
		base.Step("Удаляем из избранного объявление с несуществующим id " + nonExistentAdId)
		body := favManager.RemoveFavorite(s.T(), s.buyer.Token, nonExistentAdId, http.StatusNotFound)

		assertions.CheckError(s.T(), body, messages.ErrorNotFound, messages.FavoriteNotFound)
		base.Check("Получили 404 и сообщение «" + messages.FavoriteNotFound + "»")
	})

	s.Run("Добавление удалённого объявления", func() {
		advertisement := s.createSellerAdvertisement()
		advManager.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)
		base.Precondition("Продавец удалил объявление " + advertisement.Id)

		base.Step("Добавляем удалённое объявление в избранное")
		body := favManager.AddFavorite(s.T(), s.buyer.Token, advertisement.Id, http.StatusNotFound)

		assertions.CheckError(s.T(), body, messages.ErrorNotFound, messages.AdvertisementNotFound)
		base.Check("Получили 404 и сообщение «" + messages.AdvertisementNotFound + "»")
	})
}

// Negative: id объявления в неправильном формате
func (s *TestSuite) TestFavoritesWithInvalidAdvertisementId() {
	s.Run("Добавление", func() {
		base.Step("Добавляем в избранное объявление с id «" + invalidAdId + "»")
		body := favManager.AddFavorite(s.T(), s.buyer.Token, invalidAdId, http.StatusBadRequest)

		assertions.CheckError(s.T(), body, messages.ErrorInvalidRequest, invalidUuidMessage)
		base.Check("Получили 400 и сообщение «" + invalidUuidMessage + "»")
	})

	s.Run("Удаление", func() {
		base.Step("Удаляем из избранного объявление с id «" + invalidAdId + "»")
		body := favManager.RemoveFavorite(s.T(), s.buyer.Token, invalidAdId, http.StatusBadRequest)

		assertions.CheckError(s.T(), body, messages.ErrorInvalidRequest, messages.InvalidAdvertisementId)
		base.Check("Получили 400 и сообщение «" + messages.InvalidAdvertisementId + "»")
	})
}

// Negative: удаление из избранного объявления, которого там нет
func (s *TestSuite) TestRemoveAdvertisementNotInFavorites() {
	advertisement := s.createSellerAdvertisement()
	defer advManager.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)

	base.Step("Удаляем из избранного объявление, которое туда не добавляли")
	body := favManager.RemoveFavorite(s.T(), s.buyer.Token, advertisement.Id, http.StatusNotFound)

	assertions.CheckError(s.T(), body, messages.ErrorNotFound, messages.FavoriteNotFound)
	base.Check("Получили 404 и сообщение «" + messages.FavoriteNotFound + "»")
}
