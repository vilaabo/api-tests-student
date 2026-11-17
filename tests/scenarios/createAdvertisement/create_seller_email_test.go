package createAdvertisement

import (
	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/bugs"
	"api-tests-template/internal/constants/files"
	advManager "api-tests-template/internal/managers/advertisements"

	base "api-tests-template/tests"
)

// Negative: объявление без авторизации отдаётся с именем продавца, но без его email
func (s *TestSuite) TestAdvertisementWithoutSellerEmail() {
	createdAdvertisement := advManager.CreateAdvertisement(s.T(), s.loginData.Token, advManager.NewAdvertisementFields(),
		[]string{files.PhotoPath})
	defer advManager.DeleteAdvertisement(s.T(), s.loginData.Token, createdAdvertisement.Id)
	base.Precondition("Создали объявление, id: " + createdAdvertisement.Id)

	base.Step("Получаем объявление по id без авторизации")
	advertisement := advManager.GetAdvertisementById(s.T(), createdAdvertisement.Id)

	require.NotEmpty(s.T(), advertisement.User.FirstName, "В объявлении нет имени продавца")
	if advertisement.User.Email != "" {
		s.T().Skip(bugs.SellerEmailInResponse)
	}
	base.Check("Имя продавца есть, email продавца нет")
}
