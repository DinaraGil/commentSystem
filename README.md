## Comment System

Задача: реализация системы постов и комментариев с использованием GraphQL.

Результат: сервис поддерживает создание постов, вложенные комментарии, пагинацию и подписки на новые комментарии в реальном времени.

## Возможности

### Посты
- Получение списка всех постов
- Получение конкретного поста
- Возможность отключить комментарии к посту

### Комментарии
- Неограниченная вложенность (иерархическая структура)
- Ограничение длины комментария (до 2000 символов)
- Пагинация (`limit`, `offset`)
- Ограничение комментариев по уровню вложенности (`level`)
- Сортировка комментариев осуществляется по времени создания (от самого нового, к самому старому), сначала отображается корневой комментарий, потом его потомки

#### Комментарии возвращаются плоской структурой с указанием level и replyCommentId (на Хабре похожая реализация)

### Подписки (GraphQL Subscriptions)
- Подписка на новые комментарии к посту
- Асинхронная доставка через WebSocket

### Хранилище
- In-memory (по умолчанию)
- PostgreSQL
- Выбор через переменную окружения

## Архитектура проекта

```
.
├── cmd/server              # Точка входа (main)
├── graph                   # GraphQL схема и резолверы
│   ├── schema.graphqls
│   ├── schema.resolvers.go
│   └── loaders             # DataLoader (решение N+1)
├── internal
│   ├── models              # Модели данных
│   └── storage             # Реализация хранилищ (memory/postgres)
├── Dockerfile
├── docker-compose.yml
├── gqlgen.yml
````
## Ключевые решения

### Хранение комментариев
- Комментарии имеют поле `replyCommentID` для построения дерева
- Уровень вложенности хранится в `level`

### Пагинация
- Используется `limit` и `offset` для комментариев 1-го уровня
- Используется `level` для ограничения по уровню вложенности

### Производительность
- Используется DataLoader (для PostgreSQL) для решения проблемы N+1

### Подписки
- Используются каналы Go
- Подписчики получают обновления в реальном времени

## Запуск проекта

### 1. Клонирование

```bash
git clone https://github.com/DinaraGil/commentSystem.git
cd commentSystem
````
### 2. Запуск через Docker

#### In-memory режим

```bash
docker compose --env-file .env.memory --profile memory up --build
```

#### PostgreSQL режим

```bash
docker compose --env-file .env.postgres --profile postgres up --build
```
## GraphQL Playground

После запуска доступен по адресу:

```
http://localhost:8080/
```
## Примеры запросов
### Создание пользователя

```graphql
mutation {
  createUser(input: { username: "Ivan"}) {
    id
    username
  }
}
```
### Создание поста

```graphql
mutation {
  createPost(input: {
    userID: "1",
    content: "Message1",
    allowComment: true
  }) {
    id
  }
}
```
### Создание поста с запретом комментариев

```graphql
mutation {
  createPost(input: {
    userID: "1",
    content: "Message2",
    allowComment: false
  }) {
    id
  }
}
```

### Получение постов

```graphql
query {
  posts {
    id
    content
  }
}
```
### Добавление комментария на пост

```graphql
mutation {
  createComment(input: {
    postID: "1",
    userID: "1",
    content: "comment to post"
  }) {
    id
    level
  }
}
```
### Добавление комментария на комментарий

```graphql
mutation {
  createComment(input: {
    postID: "1",
    userID: "1",
    replyCommentID: "1"  
    content: "comment to comment"
  }) {
    id
    level
  }
}
```
### Получение всех комментариев по `postID`

```graphql
query {
    post(postID: "1") {
        id
        userID
        content
        allowComment
        createdAt
        comments {
            id
            level
            userID
            content
        }
    }
}
```
###  Limit, offset, level для комментариев по `postID`

```graphql
query {
    post(postID: "1") {
        id
        userID
        content
        allowComment
        createdAt
        comments (offset: 1, limit: 2, level: 2) {
            id
            level
            userID
            content
        }
    }
}
```

### Подписка на комментарии

```graphql
subscription {
  commentPublished(postID: "1") {
    id
    content
  }
}
```
## Тестирование

```bash
go test ./...
```
