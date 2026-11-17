package createAdvertisement

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	// не основной аккаунт: пока объявление создаётся, оно уже видно в «Моих объявлениях» без фото,
	// и тест myAdvertisement при параллельном запуске падает
	base.Precondition("Регистрация пользователя, от имени которого создаём объявления")
	s.loginData = auth.RegisterRandomUser(s.T(), "Author")
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}
