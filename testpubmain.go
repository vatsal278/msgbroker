package main

import (
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"log"
)

func main() {
	calls := sdk.NewController("http://localhost:9090")
	uuid, err := calls.RegisterPub("c11")
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(uuid)
	err = calls.PushMsg("{\\\"data\\\":\\\"hello world\\\"}", uuid, "c11")
	if err != nil {
		log.Print(err.Error())
		return
	}

}
