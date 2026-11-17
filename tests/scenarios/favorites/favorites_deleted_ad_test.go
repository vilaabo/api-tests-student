package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	favManager "api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// Positive: если продавец удалил объявление, оно пропадает из избранного покупателя
func (s *TestSuite) TestDeletedAdvertisementDisappearsFromFavorites() {
	advertisement := s.createSellerAdvertisement()
	advertisementDeleted := false
	defer func() {
		if !advertisementDeleted {
			advManager.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)
		}
	}()

	s.Run("Добавляем объявление в избранное", func() {
		base.Step("Покупатель добавляет объявление " + advertisement.Id + " в избранное")
		favManager.AddFavorite(s.T(), s.buyer.Token, advertisement.Id, http.StatusCreated)

		favorites := favManager.GetFavorites(s.T(), s.buyer.Token)
		assertions.ContainsFavorite(s.T(), favorites.Items, advertisement.Id)
		base.Check("Объявление есть в избранном")
	})

	s.Run("Продавец удаляет объявление", func() {
		base.Step("Продавец удаляет объявление через DELETE /advertisement")
		advManager.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)
		advertisementDeleted = true
	})

	s.Run("Объявление пропало из избранного покупателя", func() {
		base.Step("Получаем избранное покупателя")
		favorites := favManager.GetFavorites(s.T(), s.buyer.Token)

		assertions.NotContainsFavorite(s.T(), favorites.Items, advertisement.Id)
		base.Check("Удалённого объявления нет в избранном")
	})
}
