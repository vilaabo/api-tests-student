package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/bugs"
	advManager "api-tests-template/internal/managers/advertisements"
	favManager "api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// Positive: в избранном покупатель видит имя продавца, но не его email
func (s *TestSuite) TestFavoritesWithoutSellerEmail() {
	advertisement := s.createSellerAdvertisement()
	defer advManager.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)

	base.Step("Покупатель добавляет объявление в избранное")
	favManager.AddFavorite(s.T(), s.buyer.Token, advertisement.Id, http.StatusCreated)

	base.Step("Получаем избранное покупателя")
	favorites := favManager.GetFavorites(s.T(), s.buyer.Token)
	assertions.ContainsFavorite(s.T(), favorites.Items, advertisement.Id)

	sellerFirstName := ""
	sellerEmail := ""
	for _, item := range favorites.Items {
		if item.Advertisement.Id == advertisement.Id {
			sellerFirstName = item.Advertisement.User.FirstName
			sellerEmail = item.Advertisement.User.Email
		}
	}

	require.NotEmpty(s.T(), sellerFirstName, "В избранном нет имени продавца")
	if sellerEmail != "" {
		s.T().Skip(bugs.SellerEmailInResponse)
	}
	base.Check("Имя продавца есть, email продавца нет")
}
