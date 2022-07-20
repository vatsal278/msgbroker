<<<<<<< HEAD
## Rest API created for the Be blog system challenge

#### This Rest API was created using Golang, MySql and gorilla mux. It features functions like inserting and retreiving articles from database.

#### This api has used clean code principle 
#### This api is completely unit tested and all the errors have been handled
## Api Interface

You can test the api using post man, just import the [collection](https://github.com/vatsal278/be-blog-system-challenge/blob/bf641b8a01a9053d873a06691d24c7f212d3f5b6/docs/Blog%20System%20Collection.postman_collection.json)`collection into your postman app.

### Create an article
- Method: `POST`
- Path: `/articles`
- Request Body:
```
{
    "title": "Hello World",
    "content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
    "author": "John",
}
```
- Response Header: `HTTP 201`
- Response Body:
```
{
    "status": 201,
    "message": "Success",
    "data": {
      "id": <article_id>
    }
}
```
or
- Response Header: `HTTP <HTTP_CODE>`
- Response Body:
```
{
    "status": <HTTP_CODE>,
    "message": <ERROR_DESCRIPTION>,
    "data": null
}
```

### Get article by id
- Method: `GET`
- Path: `/articles/<article_id>`
- Response Header: `HTTP 200`
- Response Body:
```
{
    "status": 200,
    "message": "Success",
    "data": [
      {
        "id": <article_id>,
        "title":<article_title>,
        "content":<article_content>,
        "author":<article_author>,
      }
    ]
}
```
or
- Response Header: `HTTP <HTTP_CODE>`
- Response Body:
```
{
    "status": <HTTP_CODE>,
    "message": <ERROR_DESCRIPTION>,
    "data": null
}
```

### Get all article
- Method: `GET`
- Path: `/articles`
- Response Header: `HTTP 200`
- Response Body:
```
{
    "status": 200,
    "message": "Success",
    "data": [
      {
        "id": <article_id>,
        "title":<article_title>,
        "content":<article_content>,
        "author":<article_author>,
      },
      {
        "id": <article_id>,
        "title":<article_title>,
        "content":<article_content>,
        "author":<article_author>,
      }
    ]
}
```
or
- Response Header: `HTTP <HTTP_CODE>`
- Response Body:
```
{
    "status": <HTTP_CODE>,
    "message": <ERROR_DESCRIPTION>,
    "data": null
}
```
### In order to use this Api:

* clone the repo : `git clone https://github.com/vatsal278/be-blog-system-challenge.git`
* You need to specify the environment variables inside .env file, in my case it looks like 
```
DBUSER=root
DBPASS=pass
DBNAME=goblog
DBADDRESS=Database
DBPORT=3306
PORT=:8080
```
* start all the associated containers using command : `docker compose up`

### You can also run this api locally using below steps: 
* `docker run --rm --env MYSQL_ROOT_PASSWORD=pass --env MYSQL_DATABASE=goblog --publish 9095:3306 --name mysql -d mysql`
* You need to specify the environment variables inside .env file, in my case it looks like 
```
DBUSER=root
DBPASS=pass
DBNAME=goblog
DBADDRESS=localhost
DBPORT=9095
PORT=:8080
```
* Start the Api locally with command : `go run cmd/main.go`
### To check the code coverage
* `cd docs`
* `go tool cover -html=coverage`
=======
## Rest API created for the Be blog system challenge

#### This Rest API was created using Golang, MySql and gorilla mux. It features functions like inserting and retreiving articles from database.

#### This api has used clean code principle 
#### This api is completely unit tested and all the errors have been handled
## Api Interface

You can test the api using post man, just import the [collection](https://github.com/vatsal278/be-blog-system-challenge/blob/bf641b8a01a9053d873a06691d24c7f212d3f5b6/docs/Blog%20System%20Collection.postman_collection.json)`collection into your postman app.

### Create an article
- Method: `POST`
- Path: `/articles`
- Request Body:
```
{
    "title": "Hello World",
    "content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
    "author": "John",
}
```
- Response Header: `HTTP 201`
- Response Body:
```
{
    "status": 201,
    "message": "Success",
    "data": {
      "id": <article_id>
    }
}
```
or
- Response Header: `HTTP <HTTP_CODE>`
- Response Body:
```
{
    "status": <HTTP_CODE>,
    "message": <ERROR_DESCRIPTION>,
    "data": null
}
```

### Get article by id
- Method: `GET`
- Path: `/articles/<article_id>`
- Response Header: `HTTP 200`
- Response Body:
```
{
    "status": 200,
    "message": "Success",
    "data": [
      {
        "id": <article_id>,
        "title":<article_title>,
        "content":<article_content>,
        "author":<article_author>,
      }
    ]
}
```
or
- Response Header: `HTTP <HTTP_CODE>`
- Response Body:
```
{
    "status": <HTTP_CODE>,
    "message": <ERROR_DESCRIPTION>,
    "data": null
}
```

### Get all article
- Method: `GET`
- Path: `/articles`
- Response Header: `HTTP 200`
- Response Body:
```
{
    "status": 200,
    "message": "Success",
    "data": [
      {
        "id": <article_id>,
        "title":<article_title>,
        "content":<article_content>,
        "author":<article_author>,
      },
      {
        "id": <article_id>,
        "title":<article_title>,
        "content":<article_content>,
        "author":<article_author>,
      }
    ]
}
```
or
- Response Header: `HTTP <HTTP_CODE>`
- Response Body:
```
{
    "status": <HTTP_CODE>,
    "message": <ERROR_DESCRIPTION>,
    "data": null
}
```
### In order to use this Api:

* clone the repo : `git clone https://github.com/vatsal278/be-blog-system-challenge.git`
* You need to specify the environment variables inside .env file, in my case it looks like 
```
DBUSER=root
DBPASS=pass
DBNAME=goblog
DBADDRESS=Database
DBPORT=3306
PORT=:8080
```
* start all the associated containers using command : `docker compose up`

### You can also run this api locally using below steps: 
* `docker run --rm --env MYSQL_ROOT_PASSWORD=pass --env MYSQL_DATABASE=goblog --publish 9095:3306 --name mysql -d mysql`
* You need to specify the environment variables inside .env file, in my case it looks like 
```
DBUSER=root
DBPASS=pass
DBNAME=goblog
DBADDRESS=localhost
DBPORT=9095
PORT=:8080
```
* Start the Api locally with command : `go run cmd/main.go`
### To check the code coverage
* `cd docs`
* `go tool cover -html=coverage`
>>>>>>> af2ca2c15ec9d40f6fb4e869e81c0cb9aa78e330
