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

| <p id="№">№<p> | url                             | methods              |
|:--------------:|:--------------------------------|:---------------------|
|       1        | /users                          | GET / POST           |
|       2        | /users/login                    | POST                 |
|       3        | /users/refresh                  | POST                 |
|       4        | /user/{id}                      | GET / PATCH / DELETE |
|       5        | /user/id/subscribe              | POST / DELETE        |
|       6        | /posts                          | GET / POST           |
|       7        | /post/{id}                      | GET / POST / DELETE  |
|       8        | /post/{id}/like                 | POST / DELETE        |
|       9        | /post/{id}/dislike              | POST / DELETE        |
|       10       | /post/{id}/favorite             | POST / DELETE        |
|       11       | /post/{id}/comments             | GET / POST           |
|       12       | /post/{id}/comment/id           | GET / POST / DELETE  |
|       13       | /post/{id}/comment/{id}/like    | POST / DELETE        |
|       14       | /post/{id}/comment/{id}/dislike | POST / DELETE        |


<p id="1"></p>

## <a href="#№">Get Several User</a>
```
GET .../users?username=user&offset=0&limit=10
```
Происходит поиск по ?username и возвращает все совпадения.
Если ?offset и ?limit некорректные или отсутствуют, то, берется значение по умолчанию ?offset=0 и ?limit=10

Обязательные теги: `username`

Аутентификация не обязательна, возвращает <a href="#users">Users</a>

## <a href="#№">Create User</a>
```
POST .../users
```
Request body:
```
{
    "username": "user",
    "email": "user@email",
    "password": "password"
}
```
Аутентификация не обязательна, возвращает <a href="#User">User</a>

Обязательные поля: `username, email, password`

## <a href="#№">Login</a>
Авторизация на основе jwt
```
POST .../users/login
```
Request body:
```
{
    "username": "user",
    "email": "user@email",
    "password": "password"
}
```
Аутентификация не обязательна, возвращает <a href="#tokens">Tokens</a>

Обязательные поля: `username или email, password` 

## <a href="#№">Refresh</a>
Авторизация на основе jwt
```
POST .../users/refresh
```
Request body:
```
{
    "refresh": "something token"
}
```
Аутентификация не обязательна, возвращает <a href="#refresh">Refresh</a>

Обязательные поля: `refresh`

## <a href="#№">Get User</a>
```
GET .../user/{id}
```
Аутентификация не обязательна, возвращает <a href="#user">User</a>

## <a href="#№">Update User</a>
```
PATCH .../user/{id}
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

Возможные ошибки

## <a href="#№">Delete User</a>
```
DELETE .../user/{id}
```
Аутентификация обязательна.

## <a href="#№">Subscribe User</a>
```
POST .../user/{id}/subscribe
```
Аутентификация обязательна, возвращает <a href="#user">User</a>

## <a href="#№">Unsubscribe User</a>
```
DELETE .../user/{id}/subscribe
```
Аутентификация обязательна


## <a href="#№">Get Several Post</a>
```
GET .../posts ?author_id ?offset ?limit
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


[//]: #                                         (ФОРМАТ ОТВЕТОВ)
# <a href="#content" id=response>Формат ответов</a>

## <p id=tokens>Tokens</p>
```
{
    "links": {
        "self": ".../users/login",
    },
    "data": [{
        "type":"access",
        "token": "something token",
        "user:id": 1,
        "timelife": "something time"
    }, {
        "type":"refresh",
        "token": "something token",
        "user:id": 1,
        "timelife": "something time"
    }]
}
```

## <p id=refresh>Refresh</p>
```
{
    "links": {
        "self": ".../users/refresh"
    },
    "data": {
        "type":"refresh",
        "token": "something token",
        "user:id"
        "timelife": "something time"
    }
}
```

## <p id=user>User</p>
```
{
    "links": {
        "self": ".../user/1"
    },
    "data": {
        "type": "user"
        "id": 1,
        "attributes": {
            "username": "nermiabh",
            "userpic": "something url",
            "name": "David", 
            "surname": "Smyr",
            "born": "2020-07к21T12:09:00Z",
            "created": "2020-07-21T12:09:00Z",
            "subscriptionCount": 1,
            "subscriberCount": 1,
            "postCount": 0,
        },
        "relationships": {
            "self":{
                "subscribed": true,
                "subscriber": true,
                "own": false,
            }
        },
    }
}

```

## <p id=users>Users</p> 
```
{
    "links": {
        "self": ".../user?username=user&offset=0&limit=10",
        "next": ".../user?username=user&offset=10&limit=10",
        "previous": null,
    },
    "data": [{
        "type": "users",
        "id": 1,
        "attributes": {
            "username": "nermiabh",
            "userpic": "something url"
        }
    },{
        "type": "users",
        "id": 2,
        "attributes": {
            "username": "yura",
            "userpic": "something url"
        }
    }]
}
```

## <p id=user>Short User</p>
```
{
    "links": {
        "self": ".../user/1"
    },
    "data": {
        "type": "user"
        "id": 1,
        "relationships": {
            "self":{
                "subscribed": true,
                "subscriber": true,
                "own": false,
            }
        }
    }
}
```

## <p id=post>Post</p>
```
{
    "links": {
        "self": ".../post/1"
    },
    data: {
        "type": "post",
        "id": 1,
        "attributes": {
            "text": "something text",
            "media": ["something url", ]
            "created": "2020-07к21T12:09:00Z",
            "changed": true,
            "views": 0,
            "likes": 0,
            "dislikes": 0,
            "commentCount": 0,
        },
        "relationships": {
            "author":{
                "data": {
                    "type": "user", 
                    "id": "8",
                    "username": "david"
                    "userpic": "something url"
                }
            },
            "self": {
                "subscribed": true,
                "subscriber": true,
                "own": false,
                "liked": true,
                "disliked": false,
                "favorited"
            },
        }
    } 
}  
```

## <p id=posts>Posts</p>
```
{
    "links": {
        "self": ".../posts?author_id=0&offset=0$limit=10",
        "next": ".../posts?username=user&offset=10&limit=10",
        "previous": null,
    },
    data: [{
        "type": "post",
        "id": 1,
        "attributes": {
            "text": "something text",
            "media": ["something url", ]
            "created": "2020-07к21T12:09:00Z",
            "changed": true,
            "views": 0,
            "likes": 0,
            "dislikes": 0,
            "commentCount": 0,
        },
        "relationships": {
            "author":{
                "data": {
                    "type": "user", 
                    "id": "8",
                    "username": "david"
                    "userpic": "something url"
                }
            },
            "self": {
                "subscribed": true,
                "subscriber": true,
                "own": false,
            },
        }
    }, {
       ...
    }, 
  ]
}
```

## <p id=comments>Comments</p>
```
{   
    "links": {
        "self": ".../post/1/comments?offset=0&limit=5",
        "next": ".../post/1/comments?offset=5&limit=5",
        "previous": null,
    },
    data: [{
        type: "comments",
        "id": 1,
        "attributes": {
            "text": "something text",
            "created": "2020-07к21T12:09:00Z",
            "changed": true,
            "likes": 1,
            "dislikes": 0,
        },
        "relationships": {
            "author": {
                "id": 1,
                "username": "david"
                "userpic": "something url",
            }
            "post:id": 1,
            "parent:id": null,
            "self": {
                "liked": true,
                "disliked": false, 
                "own": true,
            },
        },
    }, {
        ...   
    }
 ]
}
```

## <p id=comment>Comment</p>
Как и comments только один

[//]: #                                         (ОШИБКИ)

# Коды статуса 


### 200 StatusOK
### 201 StatusCreated
Успешное создание ресурса(вместо 200)

### 400 StatusBadRequest
Произошла некая проблема из-за которой
сервер не смог обработать request

### 401 StatusUnauthorized
Возвращается если для url требуется аутентификация и она не прошла проверку.
```
{
    "errors":{
        "access": "something error"
    }
}
```

### 403 StatusForbidden
Возникает, если сервер понял запрос, 
но у пользователя нет прав совершить действие
```
{
    "errors":{
        "reason": "Вы не можете подписаться на себя"   
    }
}
```

### 404 StatusNotFound
Если статус по причине использования 
неправильного метода, то смотрите в Access-Control-Allow-Methods.
Если не найдет какой-то объект, то:
```
{
    "errors": {
        "object": "not found"
    }
}
```


### 422 StatusUnprocessableEntity
Данные не прошли валидацию
```
{
    "errors": {
        "title": "Не должен быть пустым",
        ...
    }
}
```


### 500 StatusInternalServerError
Возвращать ничего вроде не должно, но добавлю, 
так как мне будет легче разобраться в чем именно ошибка
```
{
    "errors": {
        "server": "something error"
    }.
}
```