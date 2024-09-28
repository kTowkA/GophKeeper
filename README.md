# Менеджер паролей Gophkeeper

## Сервер GophKeeper
 
### Описание
Сервер GophKeeper - это обычная реализация описанного в proto-файле gRPC-сервера. Было решено сделать некое универсальное средство для хранения данных. Поэтому запросы типа "сохранение данных кредитных кард" не описывалиьс и не реализовывались на стороне сервера. Сервер принимает массив байт и без каких-либо описательных шаблонов для принимаемых данных. Реализация такого функционала запросов должна быть представлена в клиентском приложениию Таким образом именно клиент отвечает за шифрование и дешифрование данных и представление данных в различных вариантах, а сервер останется универсальным решением для хранения данных.

Часть запросов предполагает авторизацию. Она была реализована через контекст и хранение и в нем токена доступа. Токен доступа можно получить при успешной авторизации. Те есть в запросах прямо не указывается параметр токен (пользователь) и клиент должен передать его в контексте при выполнении запроса.

Данные на сервере было решено разбить на группы (директории). Директории имеют уникальные наименования внутри пространства одного пользователя. Хранимые данные имеют уникальные наименования внутри одной директории. И директории, и данные имеют необязательное поле "описание" для хранения мета-информации.
 
#### Настройки сервера
 + Адрес по которому запускается сервер. Флаг запуска **-a**. Значение по умолчанию: **:3200**
 + Строка соединения с СУБД Postgres. Флаг запуска **-p**
 + Уровень логирования. Флаг запуска **-l**. Возможные значение
   + **debug**
   + **info** (по умолчанию)
   + **error**
 + Секрет. Переменная окружения **GOKEEPER_SECRET**. Для генерации и проверки токена доступа.

#### Тестирование
Имеется полноценный suite-тест приложения, unit-тесты отдельного функционала

#### Хранилище данных
Сервер поддерживает и работает с интерфейсом **Storager** и в настоящий момент использует реализацию его на Postgres.

#### Сборка сервера
##### Linux
```bash 
GOOS=linux GOARCH=amd64 go build -o build/server/gophkeeper-server-amd64-linux -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" server/cmd/gophkeeper_server.go
```
```bash 
GOOS=linux GOARCH=386 go build -o build/server/gophkeeper-server-386-linux -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" server/cmd/gophkeeper_server.go
```
##### Windows
```bash 
GOOS=windows GOARCH=amd64 go build -o build/server/gophkeeper-server-amd64.exe -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" server/cmd/gophkeeper_server.go
```
```bash 
GOOS=windows GOARCH=386 go build -o build/server/gophkeeper-server-386.exe -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" server/cmd/gophkeeper_server.go
```
##### MaxOS
```bash 
GOOS=darwin GOARCH=amd64 go build -o build/server/gophkeeper-server-amd64-darwin -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" server/cmd/gophkeeper_server.go
```
```bash 
GOOS=darwin GOARCH=arm64 go build -o build/server/gophkeeper-server-arm64-darwin -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" server/cmd/gophkeeper_server.go
```

## Клиенское приложение GophKeeper

### Описание
Разработанное приложение представляет собой программу с **Terminal User Interface (TUI)** (текстовым пользовательским интерфейсом).

> Предупреждение. TUI-приложение разрабатывалось в первый раз, и поэтому возникли такие временные задержки, очень много было непонятного. Тестирование было ручное, так как не разобрался как **bubble tea** можно оттестировать.

Приложение отвечает за реализацию взаимодейсвия с развернутым сервером хранения данных **GophKeeper**. 

Клиент представляет собой коллекцию дополненных моделей **bubble tea** для реализации TUI приложения.

Модели делятся на несколько основных категорий: 
 + модели с вводом значений
 + модели с отображением списков
 + моделиотображения результатов 
 + главная модель с пунктами меню для выбора. 
 
Клиент содержит функционал доступный как зарегистрированному, так и незарегистрированному пользователю. После успешной авторизации токен доступа сохраняется в контексте. Таким образом передается состояние между моделями. Кроме токена доступа в контексте передается имя пользователя и пароль. К этому контексту (по сути постоянному после авторизации) может добавляться название текущей директории и заголовок запрашиваемых данных.

Подробнее остановимся на хранении пароля в контексте. Так было решено сделать на первоначальном этапе для упрощения. В настройках модели используется интерфейс Crypter, для которого сделана реализация SimpleCrypter, основаная на симметричном шифровании одним ключем. Собственно в качестве ключа и выступает пароль пользователя.

Можно было бы сделать и следующую реализацию: пользователь регистрируется с логином и паролем и может хранить на сервере данные. А клиент, кроме логина/пароля для авторизации, запрашивает пароль для шифрования/дешифрования данных. И этот пароль не передается на сервер и таким образом мы исключаем атаку **"человек посередине"**.

Сервер предлагает универсальное решение хранения данных, а клиент может реализовать более гибкие форматы сохраняемых данных. Таким образом можно реализовать структуру "данные кредитной карты" и после получить из нее массив байт и зашифровать. Можно сделать так, просто из-за нехватки времени хотелось показать, что уже готово.

#### Настройки клиента
 + Адрес сервера. Флаг запуска **-a**. Значение по умолчанию: **:3200**
 + Файл для логгирования. Флаг запуска **-l**. Возможность логгирования предоставляет библиотека **bubble tea**.

#### Примеры
##### Генерация пароля 
![password_generate](/images/password_generate.gif)
##### Регистрация 
![registration](/images/registration.gif)
##### Авторизация
![login](/images/login.gif)
##### Сохранение данных
![save](/images/save.gif)
##### Просмотр данных
![load](/images/load.gif)
 

#### Сборка клиента
##### Linux
```bash 
GOOS=linux GOARCH=amd64 go build -o build/client/gophkeeper-client-amd64-linux -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" client/cmd/gophkeeper_client.go 
```
```bash 
GOOS=linux GOARCH=386 go build -o build/client/gophkeeper-client-386-linux -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" client/cmd/gophkeeper_client.go
```
##### Windows
```bash 
GOOS=windows GOARCH=amd64 go build -o build/client/gophkeeper-client-amd64.exe -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" client/cmd/gophkeeper_client.go 
```
```bash 
GOOS=windows GOARCH=386 go build -o build/client/gophkeeper-client-386.exe -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" client/cmd/gophkeeper_client.go
```
##### MaxOS
```bash 
GOOS=darwin GOARCH=amd64 go build -o build/client/gophkeeper-client-amd64-darwin -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" client/cmd/gophkeeper_client.go 
```
```bash 
GOOS=darwin GOARCH=arm64 go build -o build/client/gophkeeper-client-arm64-darwin -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" client/cmd/gophkeeper_client.go
```
