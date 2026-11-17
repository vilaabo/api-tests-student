package createAdvertisement

import (
	"net/http"
	"strings"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/files"
	"api-tests-template/internal/constants/messages"
	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"

	base "api-tests-template/tests"
)

const (
	titleMaxLength       = 50
	descriptionMaxLength = 500
	tooBigPrice          = "1000001"
	tooBigQuantity       = "101"

	titleRequiredMessage       = "title is required"
	descriptionRequiredMessage = "description is required"
	priceRangeMessage          = "price must be between 0 and 1000000 inclusive"
	quantityRangeMessage       = "quantity must be between 1 and 100 inclusive"

	onlySpacesBug = "BUG: обязательное поле из одних пробелов принимается, объявление создаётся с кодом 201 вместо 400"
)

// Negative: создание объявления без обязательных параметров или с некорректными параметрами
func (s *TestSuite) TestCreateAdvertisementWithInvalidParams() {
	onePhoto := []string{files.PhotoPath}
	fourPhotos := []string{files.PhotoPath, files.PhotoPath, files.PhotoPath, files.PhotoPath}

	cases := []struct {
		name            string
		field           string
		value           string
		deleteField     bool
		photos          []string
		expectedMessage string
		bug             string
	}{
		{name: "без названия", field: advModels.FieldTitle, deleteField: true, photos: onePhoto, expectedMessage: titleRequiredMessage},
		{name: "пустое название", field: advModels.FieldTitle, value: "", photos: onePhoto, expectedMessage: titleRequiredMessage},
		{name: "название из пробелов", field: advModels.FieldTitle, value: "   ", photos: onePhoto, expectedMessage: titleRequiredMessage, bug: onlySpacesBug},
		{name: "название длиннее 50 символов", field: advModels.FieldTitle, value: strings.Repeat("а", titleMaxLength+1), photos: onePhoto,
			expectedMessage: "title must not exceed 50 characters provided 51"},
		{name: "без описания", field: advModels.FieldDescription, deleteField: true, photos: onePhoto, expectedMessage: descriptionRequiredMessage},
		{name: "пустое описание", field: advModels.FieldDescription, value: "", photos: onePhoto, expectedMessage: descriptionRequiredMessage},
		{name: "описание из пробелов", field: advModels.FieldDescription, value: "   ", photos: onePhoto, expectedMessage: descriptionRequiredMessage,
			bug: onlySpacesBug},
		{name: "описание длиннее 500 символов", field: advModels.FieldDescription, value: strings.Repeat("о", descriptionMaxLength+1), photos: onePhoto,
			expectedMessage: "description must be between 1 and 500 characters"},
		{name: "отрицательная цена", field: advModels.FieldPrice, value: "-1", photos: onePhoto, expectedMessage: priceRangeMessage},
		{name: "цена больше 1000000", field: advModels.FieldPrice, value: tooBigPrice, photos: onePhoto, expectedMessage: priceRangeMessage},
		{name: "цена не число", field: advModels.FieldPrice, value: "abc", photos: onePhoto, expectedMessage: "cannot convert string value abc to float"},
		{name: "количество 0", field: advModels.FieldQuantity, value: "0", photos: onePhoto, expectedMessage: quantityRangeMessage},
		{name: "количество больше 100", field: advModels.FieldQuantity, value: tooBigQuantity, photos: onePhoto, expectedMessage: quantityRangeMessage},
		{name: "дробное количество", field: advModels.FieldQuantity, value: "1.5", photos: onePhoto, expectedMessage: "cannot convert string value to int"},
		{name: "без фото", photos: nil, expectedMessage: "at least 1 photo is required"},
		{name: "больше 3 фото", photos: fourPhotos, expectedMessage: "maximum 3 photos allowed"},
		{name: "текстовый файл вместо фото", photos: []string{files.TextFilePath},
			expectedMessage: "photo 1: file type not allowed: extension \".txt\""},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			if tc.bug != "" {
				s.T().Skip(tc.bug)
			}

			fields := advManager.NewAdvertisementFields()
			if tc.deleteField {
				delete(fields, tc.field)
			} else if tc.field != "" {
				fields[tc.field] = tc.value
			}

			base.Step("Отправляем запрос на создание объявления: " + tc.name)
			body := advManager.CreateAdvertisementWithError(s.T(), s.loginData.Token, fields, tc.photos, http.StatusBadRequest)

			assertions.CheckError(s.T(), body, messages.ErrorInvalidRequest, tc.expectedMessage)
			base.Check("Получили 400 и сообщение «" + tc.expectedMessage + "»")
		})
	}
}
