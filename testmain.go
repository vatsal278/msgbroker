package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	RSA "github.com/vatsal278/msgbroker/pkg/crypt"
)

func main() {
	var privateKey *rsa.PrivateKey
	body, err := ioutil.ReadFile("test.txt")
	err1 := json.Unmarshal(body, privateKey)
	if err != nil {
		log.Print(err1.Error())
	}
	if err != nil {

		log.Printf("failed reading data from file: %s", err)

		_, err := os.Create("privatekey.json")

		if err != nil {
			log.Printf("failed creating file: %s", err)
		}
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Print()
		}
		x, err := json.Marshal(privateKey)
		if err != nil {
			log.Printf(err.Error())
		}

		err = ioutil.WriteFile("privatekey.json", x, 0644)
		if err != nil {
			log.Fatalf("failed writing to file: %s", err)
		}
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
		res, err := RSA.RsaOaepDecrypt(string(body), *privateKey)

		log.Printf("%+v", res)
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
