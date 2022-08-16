package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	RSA "github.com/vatsal278/msgbroker/pkg/crypt"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Print()
	}

	publicKey := privateKey.PublicKey
	pubKey := RSA.KeyAsPEMStr(&publicKey)
	log.Printf("This is public key \n%v", pubKey)
	r := gin.Default()
	r.POST("/ping", func(c *gin.Context) {
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
	r.POST("/pingWoEncrypt", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Print(err)
			return
		}

		defer c.Request.Body.Close()

		var y map[string]interface{}
		err = json.Unmarshal(body, &y)
		if err != nil {
			log.Print(err)
			return
		}

		log.Printf("%+v", y)
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
	r.Run(":8086")
}

//b88847c0-a1ad-450e-a033-568e3a6cd4bc
//686f45e8-746b-4daf-bce3-c4635d90c0db
