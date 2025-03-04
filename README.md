# go-calc2

## Проект
Проект состоит из 2 элементов:
* Orchestrator - сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения.
* Agent - вычислитель, который может получить от "оркестратора" задачу, выполнить его и вернуть серверу результат.

### Принцип работы
Пользователь задаёт выражение, программа его принимает, orchestrator пердаёт ему нужный вид. Агенты через GET запрос запрашивают выражения для вычисления, если такие есть, то решает и отправляет обратно в orchestrator через POST запрос.

![Без имени](https://github.com/user-attachments/assets/898f6bb0-2385-420f-b0a6-2d0f9f0efb8f)

### Добавление проекта
```
git clone https://github.com/IDK536/go-calc2.git
```

##Запуск проекта:
Для запуска проекта выполните следующую команду:

```
go run cmd/server/main.go
```

## Эндпоинты

### Добавление выражения для вычисления:

```
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2"
}'
```

### Получение списка выражений:

```
curl --location 'localhost/api/v1/expressions'
```

### Получение выражения по ID:

```
curl --location 'localhost/api/v1/expressions/{id}'
```

## Переменные окружения для времени выполнения операций:
* TIME_ADDITION_MS: Время выполнения операции сложения в миллисекундах.
* TIME_SUBTRACTION_MS: Время выполнения операции вычитания в миллисекундах.
* TIME_MULTIPLICATIONS_MS: Время выполнения операции умножения в миллисекундах.
* TIME_DIVISIONS_MS: Время выполнения операции деления в миллисекундах.
