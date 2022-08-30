package main

import (
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"log"
)

func main() {
	//Create a new controller and pass the message broker service url to the controller.
	controller := sdk.NewMsgBrokerSvc("http://localhost:9090")
	//Call the RegisterPub function which takes in as argument the channel name for which
	//you want to subscribe and returns the uuid and error.
	//This uuid is unique and will be used for pushing the messages to subscriber.
	uuid, err := controller.RegisterPub("c11")
	//Handle the error
	if err != nil {
		log.Print(err.Error())
		return
	}
	//Log the uuid and save it for pushing messages to the channel.
	log.Print(uuid)
}
