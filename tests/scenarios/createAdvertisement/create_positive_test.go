package createAdvertisement

import (
	"fmt"
	"strconv"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/files"
	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	adTitlePrefix = "Тест создания "
	adDescription = "Объявление со всеми параметрами от автотеста"
	adPrice       = "25000.5"
	adQuantity    = "7"
)

// Positive: создаём объявление со всеми параметрами и проверяем, что оно сохранилось
func (s *TestSuite) TestCreateAdvertisementWithAllParams() {
	title := adTitlePrefix + utils.RandomString(8)
	fields := map[string]string{
		advModels.FieldTitle:       title,
		advModels.FieldDescription: adDescription,
		advModels.FieldPrice:       adPrice,
		advModels.FieldQuantity:    adQuantity,
	}
	photos := []string{files.PhotoPath, files.JpegPhotoPath, files.PhotoPath}

	base.Step("Создаём объявление «" + title + "» с тремя фото (PNG и JPEG)")
	createdAdvertisement := advManager.CreateAdvertisement(s.T(), s.loginData.Token, fields, photos)
	defer advManager.DeleteAdvertisement(s.T(), s.loginData.Token, createdAdvertisement.Id)
	base.Check("Объявление создано, id: " + createdAdvertisement.Id)

	s.Run("Объявление доступно по id и параметры совпадают с переданными", func() {
		base.Step("Получаем объявление по id")
		advertisement := advManager.GetAdvertisementById(s.T(), createdAdvertisement.Id)

		require.Equal(s.T(), createdAdvertisement.Id, advertisement.Id, "id не совпадает")
		require.Equal(s.T(), s.loginData.User.Id, advertisement.UserId, "user_id не совпадает")
		require.Equal(s.T(), title, advertisement.Title, "Название не совпадает")
		require.Equal(s.T(), adDescription, advertisement.Description, "Описание не совпадает")
		require.Equal(s.T(), adPrice, fmt.Sprint(advertisement.Price), "Цена не совпадает")
		require.Equal(s.T(), adQuantity, strconv.Itoa(advertisement.Quantity), "Количество не совпадает")
		require.Len(s.T(), advertisement.Photos, len(photos), "Количество фото не совпадает")
		base.Check("Параметры объявления совпадают с переданными при создании")
	})

	s.Run("Объявление находится поиском по полному названию", func() {
		base.Step("Ищем объявления по названию «" + title + "»")
		list := advManager.GetAdvertisementsList(s.T(), s.loginData.Token, title)

		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, title)
		base.Check("Объявление найдено в поиске")
	})

	s.Run("Фотографии загружены и доступны", func() {
		base.Step("Получаем фотографии объявления")
		advertisementPhotos := advManager.GetAdvertisementPhotos(s.T(), s.loginData.Token, createdAdvertisement.Id)

		require.Len(s.T(), advertisementPhotos, len(photos), "Количество фото не совпадает")
		for _, photo := range advertisementPhotos {
			require.NotEmpty(s.T(), photo.Id, "У фото пустой id")
			require.NotEmpty(s.T(), photo.Url, "У фото пустой url")
		}
		base.Check("Все фото на месте, у каждого есть url")
	})
}
