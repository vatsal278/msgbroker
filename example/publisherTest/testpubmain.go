package main

import (
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"log"
)

func main() {
	//Create a new controller and pass the message broker service url to the controller.
	controller := sdk.NewMsgBrokerSvc("http://localhost:9090")
	//Store the uuid which was received when registerd as publisher into uuid variable.
	uuid := "4685118a-79be-416a-bb75-f47994737b8c"
	//Call the PushMsg function which takes in as argument:
	//1. message to be pushed as a raw string,
	//2. Registered uuid which was returned after registering as publisher.
	//3. Channel name to which the message is to be pushed as string.
	//Push message function only returns error.
	err := controller.PushMsg(`{"data":"hello world"}`, uuid, "c11")
	// Handle the error
	if err != nil {
		log.Print(err.Error())
		return
	}
}
