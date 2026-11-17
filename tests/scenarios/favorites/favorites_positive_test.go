package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	favManager "api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// Positive: добавляем объявление в избранное, а потом удаляем его оттуда
func (s *TestSuite) TestAddAndRemoveFavorite() {
	advertisement := s.createSellerAdvertisement()
	defer advManager.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)

	s.Run("Добавляем объявление в избранное", func() {
		base.Step("Покупатель добавляет объявление " + advertisement.Id + " в избранное")
		favManager.AddFavorite(s.T(), s.buyer.Token, advertisement.Id, http.StatusCreated)
	})

	s.Run("Объявление появилось в избранном", func() {
		base.Step("Получаем избранное покупателя")
		favorites := favManager.GetFavorites(s.T(), s.buyer.Token)

		assertions.ContainsFavorite(s.T(), favorites.Items, advertisement.Id)
		base.Check("Объявление есть в избранном")
	})

	s.Run("Удаляем объявление из избранного", func() {
		base.Step("Покупатель удаляет объявление из избранного")
		favManager.RemoveFavorite(s.T(), s.buyer.Token, advertisement.Id, http.StatusNoContent)
	})

	s.Run("Объявление пропало из избранного", func() {
		base.Step("Получаем избранное покупателя")
		favorites := favManager.GetFavorites(s.T(), s.buyer.Token)

		assertions.NotContainsFavorite(s.T(), favorites.Items, advertisement.Id)
		base.Check("Объявления больше нет в избранном")
	})
}
