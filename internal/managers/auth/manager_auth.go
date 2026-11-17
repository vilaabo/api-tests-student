package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// Register регистрирует нового пользователя и возвращает токен с данными пользователя, ожидает 201
func Register(t *testing.T, email, password, firstName string) models.LoginOkResponse {
	request := models.RegisterRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
	}

	registerRawResponse := auth.HttpPostAuthRegister(t, request)
	defer registerRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(registerRawResponse.Body)
	require.NoError(t, errResponse)

	require.Equalf(t, http.StatusCreated, registerRawResponse.StatusCode,
		"Не удалось зарегистрировать пользователя, ответ сервера: %s", string(respArray))

	registerResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &registerResponse)
	require.NoError(t, err)

	return registerResponse
}

// RegisterRandomUser регистрирует пользователя со случайными email и паролем
func RegisterRandomUser(t *testing.T, firstName string) models.LoginOkResponse {
	email := strings.ToLower(firstName+"_"+utils.RandomString(10)) + "@test.ru"
	password := "Aa1" + utils.RandomString(10)
	return Register(t, email, password, firstName)
}
