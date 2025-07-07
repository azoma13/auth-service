# auth-service

## Описание проекта:
Сервис, который реализует часть функционала простейшей аутентификации. 

## Стек:
Язык сервиса: Go;
База данных: PostgreSQL;
Деплоя зависимостей и самого сервиса: Docker, Docker-Compose.

## Сервис предоставляет следующие конечные точки API:
- регистрация пользователя;
- аутентификация пользователя. Получение пары токенов(access и refresh) реализовано в двух вариантах: 
    - для пользователя с идентификатором (GUID) указанным в параметре запроса;
    - для пользователя с username и password указанный в теле запроса;
- обновление пары токенов;
- на получение GUID текущего пользователя;
- на деавторизацию пользователя.
Все роуты помимо регистрации и аутентификация защищены и недоступны без аутентификация в сервисе.

## Примеры
Примеры возможных запросов:
- [Регистрация](#sign-up)
- [Аутентификация](#sign-in)
- [Обновление пары токенов](#refresh)
- [Получение GUID текущего пользователя](#guid)
- [Деавторизация пользователя](#sign-out)

### Регистрация <a name="sign-up"></a>
Запрос:
```
curl --location 'http://localhost:8080/auth/sign-up' \
--header 'Content-Type: application/json' \
--data '{
    "username": "example",
    "password": "Example123!"
}'
```
Ответ: `201 Created`
```json
{
    "id":"730846ae-b356-4281-8f6c-2eb7783c9120"
}
```

### Аутентификация <a name="sign-in"></a>
Запрос для пользователя с идентификатором (GUID):
```
curl --location --request POST 'http://localhost:8080/auth/log-in?id=730846ae-b356-4281-8f6c-2eb7783c9120' \
--header 'X-Forwarded-For: 123.132.223.0'
```
Ответ: `204 No Content`
Запрос для пользователя с username и password:
```
curl --location 'http://localhost:8080/auth/sign-in' \
--header 'X-Forwarded-For: 255.255.0.0' \
--header 'Content-Type: application/json' \
--data '{
    "username": "azamatwe1",
    "password": "Qwert123!weq"
}'
```
Ответ: `204 No Content`

### Обновление пары токенов <a name="refresh"></a>
Запрос:
```
curl --location --request PUT 'http://localhost:8080/api/v1/accounts/refresh' \
--header 'X-Forwarded-For: 255.255.0.0' \
--header 'Cookie: accessToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE4NDI3NTIsImlhdCI6MTc1MTg0MDk1Miwic3ViIjoiYWNjZXNzX3Rva2VuIiwiVXNlcklkIjoiNzM0Yzk2YzgtMmJkNi00YzJmLWEzMDktMjM5NDNiMzJhZTRjIn0.vOag3CYkcRoefPPU5AjtMNFZZhK56Oxm1mOwreBM7K4; refreshToken=ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM05URTRORFk1TlRJc0ltbGhkQ0k2TVRjMU1UZzBNRGsxTWl3aWMzVmlJam9pY21WbWNtVnphRjkwYjJ0bGJpSXNJbFZ6WlhKSlpDSTZJamN6TkdNNU5tTTRMVEppWkRZdE5HTXlaaTFoTXpBNUxUSXpPVFF6WWpNeVlXVTBZeUo5LjYzN21BMlRCNW1LeXVzaVB0VmtBdTJITUdvbVdBVGRuQ1c5VE9LeGU0N1k='
```
Ответ: `204 No Content`
Запрос при изменении User-Agent:
```
curl --location --request PUT 'http://localhost:8080/api/v1/accounts/refresh' \
--header 'User-Agent: sa' \
--header 'X-Forwarded-For: 255.255.0.0' \
--header 'Cookie: accessToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE4OTY2NzIsImlhdCI6MTc1MTg5NDg3Miwic3ViIjoiYWNjZXNzX3Rva2VuIiwiVXNlcklkIjoiNzMwODQ2YWUtYjM1Ni00MjgxLThmNmMtMmViNzc4M2M5MTIwIn0.r3NViEjwAeTpBUaUFbemC6iMfNc5hhWEOXEYg4_eX7g; refreshToken=ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM05URTVNREE0TnpJc0ltbGhkQ0k2TVRjMU1UZzVORGczTWl3aWMzVmlJam9pY21WbWNtVnphRjkwYjJ0bGJpSXNJbFZ6WlhKSlpDSTZJamN6TURnME5tRmxMV0l6TlRZdE5ESTRNUzA0WmpaakxUSmxZamMzT0ROak9URXlNQ0o5LmNqWmdrTWRHdGI4dWZ5WW1YT0NuaFhpTEdWdVliMHRpU1JGWUZvZ2FKdjQ='
```
Ответ: `403 Forbidden`
```json
{
    "message": "please sign-in again"
}
```
Запрос при изменении IP:
```
curl --location --request PUT 'http://localhost:8080/api/v1/accounts/refresh' \
--header 'X-Forwarded-For: 255.255.0.1' \
--header 'Cookie: accessToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE4NDI3NjQsImlhdCI6MTc1MTg0MDk2NCwic3ViIjoiYWNjZXNzX3Rva2VuIiwiVXNlcklkIjoiNzM0Yzk2YzgtMmJkNi00YzJmLWEzMDktMjM5NDNiMzJhZTRjIn0.2HU2MteP5AKvD_jiS_xLEe1cTrEQRdp75yyO3B4Ymuw; refreshToken=ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM05URTRORFk1TmpRc0ltbGhkQ0k2TVRjMU1UZzBNRGsyTkN3aWMzVmlJam9pY21WbWNtVnphRjkwYjJ0bGJpSXNJbFZ6WlhKSlpDSTZJamN6TkdNNU5tTTRMVEppWkRZdE5HTXlaaTFoTXpBNUxUSXpPVFF6WWpNeVlXVTBZeUo5Llota01PUHBJU0lISjdLTXNrN0J0cUI2WkpkaTRubGY5enlCMkNzMFBnZUU='
```
Ответ: `204 No Content`

### Получение GUID текущего пользователя <a name="guid"></a>
Запрос:
```
curl --location 'http://localhost:8080/api/v1/accounts/guid' \
--header 'X-Forwarded-For: 255.255.0.0' \
--header 'Cookie: accessToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE4NDI3NjQsImlhdCI6MTc1MTg0MDk2NCwic3ViIjoiYWNjZXNzX3Rva2VuIiwiVXNlcklkIjoiNzM0Yzk2YzgtMmJkNi00YzJmLWEzMDktMjM5NDNiMzJhZTRjIn0.2HU2MteP5AKvD_jiS_xLEe1cTrEQRdp75yyO3B4Ymuw; refreshToken=ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM05URTRORFk1TmpRc0ltbGhkQ0k2TVRjMU1UZzBNRGsyTkN3aWMzVmlJam9pY21WbWNtVnphRjkwYjJ0bGJpSXNJbFZ6WlhKSlpDSTZJamN6TkdNNU5tTTRMVEppWkRZdE5HTXlaaTFoTXpBNUxUSXpPVFF6WWpNeVlXVTBZeUo5Llota01PUHBJU0lISjdLTXNrN0J0cUI2WkpkaTRubGY5enlCMkNzMFBnZUU='
```
Ответ: `200 OK`
```json
{
    "guid": "730846ae-b356-4281-8f6c-2eb7783c9120"
}
```

### Получение GUID текущего пользователя <a name="sign-out"></a>
Запрос:
```
curl --location --request DELETE 'http://localhost:8080/api/v1/accounts/sign-out' \
--header 'Cookie: accessToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE4OTY2NzIsImlhdCI6MTc1MTg5NDg3Miwic3ViIjoiYWNjZXNzX3Rva2VuIiwiVXNlcklkIjoiNzMwODQ2YWUtYjM1Ni00MjgxLThmNmMtMmViNzc4M2M5MTIwIn0.r3NViEjwAeTpBUaUFbemC6iMfNc5hhWEOXEYg4_eX7g; refreshToken=ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM05URTVNREE0TnpJc0ltbGhkQ0k2TVRjMU1UZzVORGczTWl3aWMzVmlJam9pY21WbWNtVnphRjkwYjJ0bGJpSXNJbFZ6WlhKSlpDSTZJamN6TURnME5tRmxMV0l6TlRZdE5ESTRNUzA0WmpaakxUSmxZamMzT0ROak9URXlNQ0o5LmNqWmdrTWRHdGI4dWZ5WW1YT0NuaFhpTEdWdVliMHRpU1JGWUZvZ2FKdjQ='
```
Ответ: `204 No Content`

## Инструкция по запуску проекта через Docker
1. Сервис собирается и запускается командой:
`docker-compose -f docker-compose.yml up -d`

