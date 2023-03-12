# Social-app документация

# <p id=content>Содержание</p>
### <a href="#routes">Routes</a>
### <a href="#response">Формат ответов</a>

## <a href=#content id=routes>Routes</a>:
Любой url api части должен начинаться с:
#### | /api/admin  | 
#### или 
#### | /api/public |
Далее urls такие:

|  <p id="№">№<p>   | url                             | methods              |
|:-----------------:|:--------------------------------|:---------------------|
| <a href="#1">1</a> | /users                          | GET / POST           |
|         2         | /users/login                    | POST                 |
|         3         | /users/refresh                  | POST                 |
|         4         | /user/{id}                      | GET / PATCH / DELETE |
|         5         | /user/id/subscribe              | POST / DELETE        |
|         6         | /posts                          | GET / POST           |
|         7         | /post/{id}                      | GET / POST / DELETE  |
|         8         | /post/{id}/like                 | POST / DELETE        |
|         9         | /post/{id}/dislike              | POST / DELETE        |
|        10         | /post/{id}/favorite             | POST / DELETE        |
|        11         | /post/{id}/comments             | GET / POST           |
|        12         | /post/{id}/comment/id           | GET / POST / DELETE  |
|        13         | /post/{id}/comment/{id}/like    | POST / DELETE        |
|        14         | /post/{id}/comment/{id}/dislike | POST / DELETE        |


<p id="1"></p>

## <a href="#№">Get Several User</a>
```
GET /users ?username=user ?offset=0 ?limit=10
```
Происходит поиск по ?username и возвращает все совпадения.
Если ?offset и ?limit некорректные, то, берется значение по умолчанию ?offset=0 и ?limit=10

Обязательные теги: `username`

Аутентификация не обязательна, возвращает <a href="#users">Users</a>

## <a href="#№">Create User</a>
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
Аутентификация не обязательна, возвращает <a href="#user">User</a>

Обязательные поля: `username, email, password`

## <a href="#№">Login</a>
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

## <a href="#№">Refresh</a>
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

## <a href="#№">Get User</a>
```
GET /user/{id}
```
Аутентификация не обязательна, возвращает <a href="#user">User</a>

## <a href="#№">Update User</a>
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
Аутентификация не обязательна, возвращает <a href="#user">User</a>

Обязательные поля: НЕТ

## <a href="#№">Delete User</a>
```
DELETE /user/{id}
```
Аутентификация обязательна, возвращает <a href="#user">User</a>

## <a href="#№">Subscribe User</a>
```
POST /user/{id}/subscribe
```
Аутентификация обязательна, возвращает <a href="#user">User</a>

## <a href="#№">Unsubscribe User</a>
```
DELETE /user/{id}/subscribe
```
Аутентификация обязательна


## <a href="#№">Get Several Post</a>
```
GET /posts ?author_id ?offset ?limit
```
Происходит поиск постов с полем ?author_id. 
Если ?offset и ?limit некорректные, то, берется значение по умолчанию ?offset=0 и ?limit=10

Аутентификация необязательна, возвращает Posts

Обязательные теги: `?author_id(временно)`

## <a href="#№">Create Post</a>
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

## <a href="#№">Get Post</a>
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

## <a href="#№">Update Post</a>
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

## <a href="#№">Delete Post</a>
```
DELETE /post/{id}
```
Аутентификация обязательна, возвращает Post

## <a href="#№">Like Post</a>
```
POST /post/{id}/like
```
Если ранее User поставил dislike, то он удаляется.

Аутентификация обязательна, возвращает Post

## <a href="#№">Unlike Post</a>
```
DELETE /post/{id}/like
```
Аутентификация обязательна, возвращает Post

## <a href="#№">Dislike Post</a>
```
POST /post/{id}/dislike
```
Если ранее User поставил like, то он удаляется.

Аутентификация обязательна, возвращает Post

## <a href="#№">Undislike Post</a>
```
DELETE /post/{id}/dislike
```
Аутентификация обязательна, возвращает Post

## <a href="#№">Favorite Post</a>
```
POST /post/{id}/favorite
```

Аутентификация обязательна, возвращает Post

## <a href="#№">Unfavorite Post</a>
```
DELETE /post/{id}/unfavorite
```
Аутентификация обязательна, возвращает Post


## <a href="#№">Get Several Comments</a>
```
GET /post/{id}/comments ?offset ?limit
```
Если ?offset и ?limit некорректные, то, берется значение по умолчанию ?offset=0 и ?limit=5

Аутентификация не обязательна, возвращает Post

Обязательные теги: ```НЕТ```


## <a href="#№">Create Comments</a>
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

## <a href="#№">Get Comment</a>
```
GET /post/{id}/comment/{id}
```
Аутентификация не обязательна, возвращает Post


## <a href="#№">Update Comment</a>
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

## <a href="#№">Delete Comment</a>
```
DELETE /post/{id}/comment/{id}
```
Аутентификация не обязательна


## <a href="#№">Like Comment</a>
```
POST /post/{id}/comment/{id}/like
```
Если ранее User поставил dislike, то он удаляется.

Аутентификация обязательна, возвращает Comment

## <a href="#№">Unlike Comment</a>
```
DELETE /post/{id}/comment/{id}/like
```
Аутентификация обязательна, возвращает Comment

## <a href="#№">Dislike Comment</a>
```
POST /post/{id}/comment/{id}/dislike
```
Если ранее User поставил like, то он удаляется.

Аутентификация обязательна, возвращает Comment

## <a href="#№">Undislike Comment</a>
```
DELETE /post/{id}/comment/{id}/dislike
```
Аутентификация обязательна, возвращает Comment


# <a href="#content" id=response>Формат ответов</a>
### Предисловие:
Некоторые ниже перечисленные поля ответов могут 
вовсе не выводиться подразумевая, 
что в нем хранится значение по умолчанию
(для поля с bool типом это false, для числа - 0, для строки - "") или ничего не хранится
PS. Могу изменить если это неудобно для клиента, в начале я думал что это фифа, а оказалось не совсем

## <p id=user>User</p>
```
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
```

## <p id=users>Users</p> 
```
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
```

## <p id=post>Post</p>
```
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
        "is_own": true
    }    
```

## <p id=posts>Posts</p>
```
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
```

## <p id=comments>Comments</p>
```
comments: [{
        "id": 1,
        "author_id": 1, 
        "author_userpic": "something url",
        "post_id": 1,
        "parent_id": null,
        "text": "something text",
        "date_of_creation": "22-10-2022 time",
        "is_changed": true,
        "likes": 1,
        "dislikes": 0,
        "is_own": true,
        "is_liked": true,
        "is_disliked": false, 
    }, {
        ...   
    }
 ]
```

## <p id=comment>Comment</p>
Как и comments только один

