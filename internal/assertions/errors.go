package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// CheckError проверяет поля error и message в теле ответа с ошибкой
func CheckError(t *testing.T, body, expectedError, expectedMessage string) {
	require.Equal(t, expectedError, gjson.Get(body, "error").String(),
		"Некорректное значение error параметра, тело ответа: %s", body)
	require.Equal(t, expectedMessage, gjson.Get(body, "message").String(),
		"Некорректное значение message параметра, тело ответа: %s", body)
}
