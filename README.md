## This is a message broker service API.

#### This Rest API was created using Golang. It features functions like registering publisher or subscriber, publishing messages or extracting updates.

#### This api has used clean code principle 
#### This api is completely unit tested and all the errors have been handled
## Api Interface

You can test the api using post man, just import the [collection](https://github.com/vatsal278/be-blog-system-challenge/blob/bf641b8a01a9053d873a06691d24c7f212d3f5b6/docs/Blog%20System%20Collection.postman_collection.json)`collection into your postman app.

### Register a Publisher
- Method: `POST`
- Path: `/register/publisher`
- Request Body:
```
{
    "channel": "channel x",
}
```
- Response Header: `HTTP 201`
- Response Body:
```
{
    "status": 201,
    "message": "Successfully Registered as publisher to the channel",
    "data": {
      "id": "8a6d08a1-0cc3-451a-a66e-474fbf16a781"
    }
}
```

### Register a subscriber
- Method: `POST`
- Path: `/register/subscriber`
- Request Body:
```
{
    "callback":{
          "httpmethod":"POST",
          "callbackUrl":"http://localhost:8080/ping",
          "key":"",
          },
    "channel": "channel x",
}
```
- Response Header: `HTTP 201`
- Response Body:
```
{
    "status": 201,
    "message": "Successfully Registered as subscriber to the channel",
    "data": null
}
```

### Push Messages
- Method: `POST`
- Path: `/publish`
- Request Body:
```
{
    "publisher": {
        "id": "8a6d08a1-0cc3-451a-a66e-474fbf16a781",
        "channel": "c1"
    },
    "update": "{\"data\":\"hello world\"}"
}
```
- Response Header: `HTTP 200`
- Response Body:
```
{
    "status": 200,
    "message": "Notified All Subscriber",
    "data": null
}
```

### In order to use this Api:

* clone the repo : `git clone https://github.com/vatsal278/msgbroker`
* build the docker file using command : `docker build -t msgbrokersvc .`
* run the docker container : `docker run --rm --env PORT=9090 -p 9099:9090 msgbrokersvc`

### You can also run this api locally using below steps: 
* Start the MsgBrokerSvc locally with command : `go run cmd/main.go`

### In order to use the SDK functions:
* create a new controller and pass url to the msg broker server : ```controller := sdk.NewController("http://localhost:9090")```
* To register as a publisher : `uuid, err := controller.RegisterPub("channel")`
* To push messages to all the subscribers to that channel pass the message, uuid which you got when registering as publisher along with channel name : ``err := calls.PushMsg(`{"data":"hello world"}`, uuid, "channel")``
* To register as a subscriber pass the subscriber details as httpmethod, subscriberurl, publickey(optional:if want to get non encrypted message otherwise send an empty string) : `err := controller.RegisterSub("POST", "http://localhost:9091/ping", pubKey, "c11")`
* To extract the message sent by the publisher : 
* `extractMsg := controller.ExtractMsg(privateKey) `
* `s, err := extractMsg(io.ReadClosure)`

### In order to test the publisher service:
* start the msgbrokersvc
* cd publisherTest
* go run main.go

### In order to test the subscriber service:
* start the msgbrokersvc
* cd subscriberTest
* go run main.go




