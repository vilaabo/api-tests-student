package api_runner

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/steinfletcher/apitest"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"

	maxAttempts = 5
	retryDelay  = time.Second
)

// retryTransport повторяет запрос, если сервер ответил 503: так nginx ограничивает частоту запросов
type retryTransport struct{}

// RoundTrip отправляет запрос и при ответе 503 повторяет его с паузой
func (rt retryTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var body []byte
	if request.Body != nil {
		var err error
		body, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		request.Body.Close()
	}

	for attempt := 1; ; attempt++ {
		attemptRequest := request.Clone(request.Context())
		attemptRequest.Body = io.NopCloser(bytes.NewReader(body))

		response, err := http.DefaultTransport.RoundTrip(attemptRequest)
		if err != nil || response.StatusCode != http.StatusServiceUnavailable || attempt == maxAttempts {
			return response, err
		}
		response.Body.Close()

		delay := time.Duration(attempt) * retryDelay
		log.Printf("Сервер ответил 503, повторяем запрос через %v (попытка %d из %d)", delay, attempt+1, maxAttempts)
		time.Sleep(delay)
	}
}

var httpClient = &http.Client{
	Transport: retryTransport{},
	Timeout:   time.Minute,
}

type ApiTest struct {
	token     string
	userAgent string
	host      string
}

func New() *ApiTest {
	return &ApiTest{
		userAgent: defaultUserAgent,
		host:      os.Getenv("API_URL"),
	}
}

func GetRunner() *ApiTest {
	return New()
}

// Auth Установка Bearer token-а для прохождения авторизации на приватные ручки
func (at *ApiTest) Auth(token string) *ApiTest {
	at.token = token
	return at
}

func (at *ApiTest) Create() *apitest.APITest {
	apitestNew := apitest.New().EnableNetworking(httpClient)
	if os.Getenv("DEBUG") == "true" {
		apitestNew = apitestNew.Debug()
	}
	return apitestNew.
		Intercept(func(request *http.Request) {
			request.Header.Set("User-Agent", at.userAgent)

			if err := mergeServiceUrls(at.host, request.URL); err != nil {
				panic(fmt.Sprintf("Не удалось применить хост из API_URL: %s", err.Error()))
			}

			if len(at.token) != 0 {
				request.Header.Add("Authorization", "Bearer "+at.token)
			}
		})
}

func mergeServiceUrls(host string, r *url.URL) error {
	if host == "" {
		return fmt.Errorf("API_URL не задан: проверьте переменные окружения")
	}

	urlParsed, err := url.Parse(host)
	if err != nil {
		return fmt.Errorf("host cannot be parsed: %s", err.Error())
	}

	if urlParsed.Path != "" {
		r.Path = path.Join(urlParsed.Path, r.Path)
	}

	r.Scheme = urlParsed.Scheme
	r.Host = urlParsed.Host
	return nil
}
