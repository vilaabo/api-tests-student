# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1 или выше
4. Склонируйте репозиторий по HTTPS и перейдите в его папку:

```bash
   git clone https://github.com/vilaabo/api-tests-student.git
   cd api-tests-student
```

5. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

6. В файле `.env.override` (создайте его в корне, если нет) укажите логин и пароль от своего аккаунта на TestBoard:

```
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
```

7. Для запуска всех тестов выполните команду в корне проекта (`-count=1` отключает кэш, чтобы тесты
   действительно ходили в API при каждом запуске):

```bash
  go test -count=1 -v ./...
```

8. Отдельный пакет с тестами можно запустить командой:

```bash
  go test -count=1 -v ./tests/scenarios/createAdvertisement/
  go test -count=1 -v ./tests/scenarios/favorites/
```

Тест `myAdvertisement` из шаблона ожидает, что у аккаунта из `.env.override` уже есть хотя бы одно объявление с фото.
Остальным тестам заранее ничего готовить не нужно.

## Что покрыто тестами

Создание объявления, пакет `tests/scenarios/createAdvertisement`:

* `create_positive_test.go` — объявление со всеми параметрами: GET по id возвращает те же данные, объявление находится
  поиском по полному названию, фото доступны через `GET /advertisements/{id}/photos`;
* `create_negative_test.go` — без обязательных полей и с некорректными значениями (длина, цена, количество, фото);
* `create_unauthorized_test.go` — запрос без токена и с неверным токеном;
* `create_seller_email_test.go` — в объявлении нет email продавца (баг №1).

Избранное, пакет `tests/scenarios/favorites`:

* `favorites_positive_test.go` — добавили объявление в избранное и удалили из него;
* `favorites_deleted_ad_test.go` — удалённое продавцом объявление пропадает из избранного;
* `favorites_not_found_test.go` — несуществующее объявление, id в неверном формате, объявления нет в избранном;
* `favorites_unauthorized_test.go` — все ручки избранного без токена и с неверным токеном;
* `favorites_business_rules_test.go` — повторное добавление (409), своё объявление (403), запрос без
  `advertisement_id` (400), лимит 50 объявлений (400);
* `favorites_seller_email_test.go` — в избранном нет email продавца (баг №1).

## Тестовые данные

* Тесты создают объявления со случайными названиями и удаляют их в конце, поэтому их можно запускать много раз подряд.
* Своё объявление добавить в избранное нельзя, поэтому тесты избранного сами регистрируют продавца через
  `POST /auth/register`. Покупатель — пользователь из `.env.override`.
* Для теста лимита регистрируется отдельный покупатель, чтобы не заполнять избранное основного аккаунта.

## Ограничение частоты запросов

Если отправлять больше ~10 запросов в секунду, сервер начинает отвечать `503`. Тесты из разных пакетов идут
параллельно, поэтому в таком случае запрос автоматически повторяется (до 5 попыток), в логе появляется строка
`Сервер ответил 503, повторяем запрос`.

## Найденные дефекты

Тесты на эти дефекты помечены как `SKIP`, в выводе видна причина `BUG: ...`.

1. Баг №1 из баг-листа: API отдаёт email продавца. Кроме ленты, email приходит в `GET /advertisement?id=`
   (без авторизации) и в `GET /favorites` (покупатель видит email продавца).
2. Обязательные поля `title` и `description` принимают строку из одних пробелов: объявление создаётся
   с кодом `201` вместо `400`.
3. `DELETE /advertisement?id=` с несуществующим или уже удалённым id возвращает `403 You can only delete your own
   advertisements`, хотя по swagger должен быть `404`, а по ТЗ `400`.
4. `POST /favorites` без `advertisement_id` возвращает в `message` внутренний текст валидатора:
   `Key: 'AddFavoriteRequest.AdvertisementID' Error:Field validation for 'AdvertisementID' failed on the 'required' tag`.
5. `GET /advertisements/{id}/photos` требует авторизацию, но в swagger ответ `401` не описан.

## Логирование

Для включения логирования request/response в консоль: выставьте переменную окружения `DEBUG=true` в
`.env.override`

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```
