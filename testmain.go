package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/ping", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Print(err)
			return
		}
		log.Print(1)
		defer c.Request.Body.Close()
		var x map[string]interface{}
		err = json.Unmarshal(body, &x)
		log.Print(2)
		if err != nil {
			log.Print(err)
			return
		}
		log.Print(3)
		log.Printf("%+v", x)
	})
	r.GET("/pong", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Print(err)
			return
		}
		defer c.Request.Body.Close()
		var x map[string]interface{}
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Print(err)
			return
		}
		log.Printf("%+v", x)
	})
	r.Run(":8082")
}

//b88847c0-a1ad-450e-a033-568e3a6cd4bc
//686f45e8-746b-4daf-bce3-c4635d90c0db
