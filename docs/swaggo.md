# Swaggo
https://github.com/swaggo/swag

в `/cmd/backend/main.go` задаем конфиг и описание для сваггера

Для обновления/генерации спецификации:
```
swag init -g ./cmd/backend/main.go -o ./api   
```