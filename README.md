## Message Broker Service

* This Rest API was created using Golang. It features functions like registering publisher or subscriber, publishing messages or extracting updates.
* This utilises RSA encryption and the key can be passed to the functions for encryption/decryption, following a required format.
* This api has used clean code principle 
* This api is completely unit tested and all the errors have been handled

#### Formatting the `RSA key`
* `RSA key` can be supplied by the user
* The `RSA public key` follows a specific format that can be obtained by the pkg function
```
import "github.com/vatsal278/msgbroker/pkg/crypt"
.
.
pubKey := crypt.KeyAsPEMStr(RSA public key)
//pubKey is the string formating of the RSA public key that will be used for subscriber registration
```
## Starting the message broker service

* Running locally : 
```
go run cmd/main.go
```
* Running via docker :
```
docker build -t msgbrokersvc .
docker run -p 9090:9090 msgbrokersvc
```
## API Spec

You can test the api using post man, just import the [collection](./docs/New%20Collection.postman_collection.json) into your postman app.

### Register a Publisher
- Method: `POST`
- Path: `/register/publisher`
- Request Body:
```
{
    "channel": "channel",
}
```
- Response Header: `HTTP 201`
- Response Body:
```
{
    "status": 201,
    "message": "Successfully Registered as publisher to the channel",
    "data": {
      "id": "Uuid"
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
          "callbackUrl":"URL",
          "key":"RSA public  key (in pem string format)"
    },
    "channel": "channel"
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
        "id": "Uuid",
        "channel": "channel"
    },
    "update": "message"
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


### In order to use the SDK functions:
* `go get` the package 
```
go get github.com/vatsal278/msgbroker
```
* `import` the `sdk` package in the source code
```
import "github.com/vatsal278/msgbroker/pkg/sdk"
```
* Get an instance to the SDK Wrapper, Passing in the url to a running message broker service.
```
s := sdk.NewMsgBrokerSvc("Message broker service url")
```
* Register a publisher to push messages to a `channel`. A Uuid of the publisher will be returned that needs to be used to push message to the `channel`.
```
uuid, _ := s.RegisterPub("channel")
```
* Register a subscriber to receive messages on a `channel`. 
Send in the `HTTP method`, `URL` to the subscriber endpoint, optional `RSA public key` 
(this will be used by the message broker to encrypt the message before notifying the subscriber) and the `channel`
```
_ = s.RegisterSub("HTTP method", "URL", "RSA public key", "channel")
```
* Push `message` to all the subscribers registered on the `channel` using the registered publishers `Uuid`.
```
_ = s.PushMsg(`message`, Uuid, "channel")
```
* Extract `message` pushed to the subscriber registered on the `channel`, additionally decrypting with `RSA private key`,
if opted for receiving encrypted messages during subscriber registration.
`source` is of type `io.ReadClosure` which can be requests body.
```
MsgExtractor := s.ExtractMsg(RSA private key)
msg, _ := MsgExtractor(source)
```
* Examples of the sdk usage can be found [here](./example)

## Additional read
* [Docs](./docs/README.md)
* To check the code coverage 
```
cd docs
go tool cover -html=coverage
```
  



