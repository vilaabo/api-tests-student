package favorites

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/constants/files"
	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	buyer  models.LoginOkResponse
	seller models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация покупателя с кредами из переменных окружения")
	s.buyer = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Регистрация продавца: своё объявление в избранное добавить нельзя, поэтому нужен второй пользователь")
	s.seller = auth.RegisterRandomUser(s.T(), "Seller")
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// createSellerAdvertisement создаёт объявление от имени продавца
func (s *TestSuite) createSellerAdvertisement() advModels.Advertisement {
	advertisement := advManager.CreateAdvertisement(s.T(), s.seller.Token, advManager.NewAdvertisementFields(), []string{files.PhotoPath})
	base.Precondition("Продавец создал объявление, id: " + advertisement.Id)
	return advertisement
}
