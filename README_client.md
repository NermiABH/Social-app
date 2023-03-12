# Social-app документация

## Все urls:
Любой url api части должен начинаться с :

#### | api/admin  | 
#### или 
#### | api/public |


Далее urls такие:

|        №        | url                             | methods              |
|:---------------:|:--------------------------------|:---------------------|
|        1        | /users                          | GET / POST           |
|        2        | /users/login                    | POST                 |
|        3        | /users/refresh                  | POST                 |
|        4        | /user/{id}                      | GET / PATCH / DELETE |
|        5        | /user/id/subscribe              | POST / DELETE        |
|        6        | /posts                          | GET / POST           |
|        7        | /post/{id}                      | GET / POST / DELETE  |
|        8        | /post/{id}/like                 | POST / DELETE        |
|        9        | /post/{id}/favorite             | POST / DELETE        |
|       10        | /post/{id}/comments             | GET / POST           |
|       11        | /post/{id}/comment/id           | GET / POST / DELETE  |
|       12        | /post/{id}/comment/{id}/like    | POST / DELETE        |
|       13        | /post/{id}/comment/{id}/dislike | POST / DELETE        |



## Get Several User
```
GET /users ?username=user ?offset=0 ?limit=10
```
Происходит поиск по ?username и возвращает все совпадения.
Если ?offset и ?limit некорректные, то , берется значение по умолчанию ?offset=0 и ?limit=10

Обязательные теги: `username`

Аутентификация не обязательна, возвращает Users
## Create User
```
POST /users
```
Request body:
```
{
    "username": "user",
    "email": "user@email",
    "password": "password"
}
```
Аутентификация не обязательна, возвращает Tokens

Обязательные поля: `username, email, password`

## Login:
Авторизация на основе jwt
```
POST /users/login
```
Request body:
```
{
    "username": "user",
    "email": "user@email",
    "password": "password"
}
```
Аутентификация не обязательна, возвращает Tokens

Обязательные поля: `username или email, password` 

## Refresh:
Авторизация на основе jwt
```
POST /users/refresh
```
Request body:
```
{
    "refresh": "token"
}
```
Аутентификация не обязательна, возвращает Tokens

Обязательные поля: `refresh`

## Get User
```
GET /user/{id}
```
Аутентификация не обязательна, возвращает User

## Update User
```
PATCH /user/{id}
```
Request body:
```
{
    "username": "user",
    "email": "user@email",
    "userpic": "url",
    "name": "user",
    "surname": ""
}
```
Аутентификация не обязательна, возвращает User

Обязательные поля: НЕТ

## Delete User
```
DELETE /user/{id}
```
Аутентификация обязательна, возвращает User

## Subscribe User
```
POST /user/{id}/subscribe
```
Аутентификация обязательна, возвращает User

## Unsubscribe User
```
DELETE /user/{id}/subscribe
```
Аутентификация обязательна, возвращает User


## Get Several Post
```
GET /posts ?author_id ?offset ?limit
```
Происходит поиск постов с полем ?author_id. 
Если ?offset и ?limit некорректные, то, берется значение по умолчанию ?offset=0 и ?limit=10

Аутентификация необязательна, возвращает Posts

Обязательные теги: `?author_id(временно)`

## Create Post
```
POST /posts
```
Request body:
```
{
    "text": "something text",
    "object": "something url",
}
```
Аутентификация обязательна, возвращает Post

Обязательные поля: `text или object не пустые`

## Get Post
```
GET /users/{id}
```
Request body:
```
{
    "text": "something text",
    "object": "something url",
}
```
Аутентификация не обязательна, возвращает Post

Обязательные поля: `text или object не пустой`

## Update Post
```
PATCH /users/{id}
```
Request body:
```
{
    "text": "something text",
    "object": "something url",
}
```
Аутентификация обязательна, возвращает Post

Обязательные поля: `text или object не пустой`

## Delete Post
```
DELETE /post/{id}
```
Аутентификация обязательна, возвращает Post

## Like Post
```
POST /post/{id}/like
```
Если ранее User поставил dislike, то он удаляется.

Аутентификация обязательна, возвращает Post

## Unlike Post
```
DELETE /post/{id}/like
```
Аутентификация обязательна, возвращает Post

## Dislike Post
```
POST /post/{id}/dislike
```
Если ранее User поставил like, то он удаляется.

Аутентификация обязательна, возвращает Post

## Undislike Post
```
DELETE /post/{id}/dislike
```
Аутентификация обязательна, возвращает Post

## Favorite Post
```
POST /post/{id}/favorite
```

Аутентификация обязательна, возвращает Post

## Unfavorite Post
```
DELETE /post/{id}/unfavorite
```
Аутентификация обязательна, возвращает Post


## Get Several Comments 
```
GET /post/{id}/comments ?offset ?limit
```
Если ?offset и ?limit некорректные, то, берется значение по умолчанию ?offset=0 и ?limit=5

Аутентификация не обязательна, возвращает Post

Обязательные теги: ```НЕТ```


## Create Comments 
```
POST /post/{id}/comments
```
Request body:
```
{
    "text": "something text"
}
```
Аутентификация обязательна, возвращает Post

Обязательные поля: `text`

## Get Comment
```
GET /post/{id}/comment/{id}
```
Аутентификация не обязательна, возвращает Post


## Update Comment
```
POST /post/{id}/comment/{id}
```
Request body:
```
{
    "text": "something text"
}
```
Аутентификация обязательна, возвращает Post

Обязательные поля: `text`

## Delete Comment
```
DELETE /post/{id}/comment/{id}
```
Аутентификация не обязательна


## Like Comment
```
POST /post/{id}/comment/{id}/like
```
Если ранее User поставил dislike, то он удаляется.

Аутентификация обязательна, возвращает Comment

## Unlike Comment
```
DELETE /post/{id}/comment/{id}/like
```
Аутентификация обязательна, возвращает Comment

## Dislike Comment
```
POST /post/{id}/comment/{id}/dislike
```
Если ранее User поставил like, то он удаляется.

Аутентификация обязательна, возвращает Comment

## Undislike Comment
```
DELETE /post/{id}/comment/{id}/dislike
```
Аутентификация обязательна, возвращает Comment


# Формат ответов
### Предисловие:
Некоторые ниже перечисленные поля ответов могут 
вовсе не выводиться подразумевая, 
что в нем хранится значение по умолчанию
(для поля с bool типом это false, для числа - 0, для строки - "") или ничего не хранится
PS. Могу изменить если это неудобно для клиента, в начале я думал что это фифа, а оказалось не совсем

## User
```
{
    user: {
        "id": 1,
        "username": "nermiabh",
        "userpic": "something url",
        "name": "David", 
        "surname": "Smyr",
        "date_of_birth": "25-10-2003",
        "date_of_creation": "25-10-2003 time",
        "subscriptions_count": 1,
        "subscribers_count": 1,
        "posts_count": 0,
        "favorites_posts": 0,
        "is_subscription": true,
        "is_subscriber": true,
        "is_own": false,
    }
}
```

## Users 
```
{
    users: [{
        "id": 1,
        "username": "nermiabh",
        "userpic": "something url"
      },
      {
        "id": 2,
        "username": "yura",
        "userpic": "something url"
      }
    ]
}
```

## Post
```
{
    post: {
        "id": 1,
        "author_id": 1,
        "text": "something text",
        "object": "something url",
        "comments_count": 0,
        "date_of_creation": "10.12.2022 time",
        "likes": 0,
        "dislikes": 0,
        "is_liked": false,
        "is_disliked": true,
        "is_favorited": true,
    }
}
```

## Posts
```
{
    posts: [{
       "id": 1,
       "text": "something text",
       "object": "something url",
       "date_of_creation": "10-10-2022 time", 
    }, {
       "id": 2,
       "text": "something text",
       "object": "something url",
       "date_of_creation": "10-12-2022 time", 
    }, 
    ]
}
```

## Comment
```
   
```
