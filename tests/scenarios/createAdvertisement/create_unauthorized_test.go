package createAdvertisement

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/files"
	"api-tests-template/internal/constants/messages"
	advManager "api-tests-template/internal/managers/advertisements"

	base "api-tests-template/tests"
)

// Negative: создание объявления без авторизации или с неверным токеном
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	cases := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{"без токена", "", messages.AuthorizationHeaderRequired},
		{"с неверным токеном", base.InvalidToken, messages.InvalidAuthorization},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			base.Step("Отправляем запрос на создание объявления " + tc.name)
			body := advManager.CreateAdvertisementWithError(s.T(), tc.token, advManager.NewAdvertisementFields(),
				[]string{files.PhotoPath}, http.StatusUnauthorized)

			assertions.CheckError(s.T(), body, messages.ErrorUnauthorized, tc.expectedMessage)
			base.Check("Получили 401 и сообщение «" + tc.expectedMessage + "»")
		})
	}
}
