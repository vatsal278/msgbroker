package main

import (
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"log"
)

func main() {
	calls := sdk.NewController("http://localhost:9090")
	uuid := "4685118a-79be-416a-bb75-f47994737b8c"
	err := calls.PushMsg(`{"data":"hello world"}`, uuid, "c11")
	if err != nil {
		log.Print(err.Error())
		return
	}
}