package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/messages"
	advManager "api-tests-template/internal/managers/advertisements"
	favManager "api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// Negative: запросы к избранному без авторизации или с неверным токеном
func (s *TestSuite) TestFavoritesUnauthorized() {
	advertisement := s.createSellerAdvertisement()
	defer advManager.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)

	cases := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{"без токена", "", messages.AuthorizationHeaderRequired},
		{"с неверным токеном", base.InvalidToken, messages.InvalidAuthorization},
	}

	for _, tc := range cases {
		s.Run("Получение избранного "+tc.name, func() {
			base.Step("Запрашиваем избранное " + tc.name)
			body := favManager.GetFavoritesWithError(s.T(), tc.token, http.StatusUnauthorized)

			assertions.CheckError(s.T(), body, messages.ErrorUnauthorized, tc.expectedMessage)
			base.Check("Получили 401 и сообщение «" + tc.expectedMessage + "»")
		})

		s.Run("Добавление в избранное "+tc.name, func() {
			base.Step("Добавляем объявление в избранное " + tc.name)
			body := favManager.AddFavorite(s.T(), tc.token, advertisement.Id, http.StatusUnauthorized)

			assertions.CheckError(s.T(), body, messages.ErrorUnauthorized, tc.expectedMessage)
			base.Check("Получили 401 и сообщение «" + tc.expectedMessage + "»")
		})

		s.Run("Удаление из избранного "+tc.name, func() {
			base.Step("Удаляем объявление из избранного " + tc.name)
			body := favManager.RemoveFavorite(s.T(), tc.token, advertisement.Id, http.StatusUnauthorized)

			assertions.CheckError(s.T(), body, messages.ErrorUnauthorized, tc.expectedMessage)
			base.Check("Получили 401 и сообщение «" + tc.expectedMessage + "»")
		})
	}
}
