# matcha
API для социальной сети знакомств  

## Оглавление
- [Конфигурация](#configuration)
- [Развертывание](#deployment)

<a name="configuration"></a>
## Конфигурация
Конфигурация происходит следующим образом:  
1. Читается конфиг по пути, указанному в переменной окружения `CONFIG_PATH`.
Файл должен иметь формат YAML и иметь определенную структуру.
2. Читаются оставшиеся настройки из переменных окружения.
При дублировании настроек переменные окружения затирают параметры конфига.

Список переменных окружения:  
```
CONFIG_PATH=configs/config.yaml
LOGGER_LEVEL=debug
LOGGER_FORMAT=default
PG_ADDRESS=0.0.0.0:5432
PG_USER=matcha
PG_PASSWORD=123
PG_DATABASE=matcha
REDIS_ADDRESS=0.0.0.0:6379
JWT_SIGNING_KEY=some_key
COOKIE_HASH_KEY=some_key
COOKIE_BLOCK_KEY=some_key
COOKIE_DOMAIN=matcha.com
EMAIL_SERVER_ADDRESS=smtp.gmail.com:587
EMAIL_USERNAME=user@test.com
EMAIL_PASSWORD=some_password
MINIO_ENDPOINT=0.0.0.0:9000
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=minio123
```

<a name="deployment"></a>
## Развертывание
1. Для того, чтобы развернуть сервис в docker:  
```docker-compose up -d```  

    Опустить контейнеры:  
```docker-compose down```  
2. Чтобы выполнить начальную миграцию для базы данных, нужно установить <a href="https://github.com/golang-migrate/migrate">эту утилиту</a> и выполнить команду:  
```migrate -path ./schema -database 'postgres://matcha:123@localhost:54320/matcha?sslmode=disable' up```  

    Откатить миграцию:  
```migrate -path ./schema -database 'postgres://matcha:123@localhost:54320/matcha?sslmode=disable' down```  
