package advertisementSearch

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/files"
	advList "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	existingAdTitle = "Ехал Грека Через Реку" // полное название объявления, которое создаёт сьют
	searchBeginning = "Ехал"                  // подстрока из начала
	searchMiddle    = "Грека Через"           // подстрока из середины
	searchEnd       = "Реку"                  // подстрока из конца
)

type TestSuite struct {
	suite.Suite
	loginData     models.LoginOkResponse
	advertisement advModels.Advertisement
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Проверяем, что константы заданы")
	require.NotEmpty(s.T(), existingAdTitle, "задайте existingAdTitle")
	require.NotEmpty(s.T(), searchBeginning, "задайте searchBeginning")
	require.NotEmpty(s.T(), searchMiddle, "задайте searchMiddle")
	require.NotEmpty(s.T(), searchEnd, "задайте searchEnd")

	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Создаём объявление «" + existingAdTitle + "», которое будем искать")
	fields := advList.NewAdvertisementFields()
	fields[advModels.FieldTitle] = existingAdTitle
	s.advertisement = advList.CreateAdvertisement(s.T(), s.loginData.Token, fields, []string{files.PhotoPath})
}

func (s *TestSuite) TearDownSuite() {
	base.Precondition("Удаляем объявление, созданное для поиска")
	advList.DeleteAdvertisement(s.T(), s.loginData.Token, s.advertisement.Id)

	base.TearDownSuite()
}

// Positive: ищем по подстроке из середины названия и находим своё объявление
func (s *TestSuite) TestSearchByTitle() {
	list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, searchMiddle)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)
	})
	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
	})
}

// Positive: подстрока в начале, середине, конце названия и в разных регистрах
func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", utils.MixCase(searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	}
}
