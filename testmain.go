package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vatsal278/msgbroker/pkg/crypt"
	"io/ioutil"
	"log"
)

func main() {
	var privateKey *rsa.PrivateKey
	body, err := ioutil.ReadFile("privatekey.pem")

	if err != nil {

		log.Printf("failed reading data from file: %s", err)

		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Print(err.Error())
			return
		}
		KeyPem := string(pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			},
		))

		err = ioutil.WriteFile("privatekey.pem", []byte(KeyPem), 0644)
		if err != nil {
			log.Printf("failed writing to file: %s", err)
			return
		}
		log.Print("succesfully saved PEM Key to file")
	} else {
		spkiBlock, _ := pem.Decode(body)
		if spkiBlock == nil || spkiBlock.Type != "RSA PRIVATE KEY" {
			err := errors.New("failed to decode PEM block containing public key")
			log.Print(err.Error())
			return
		}
		privateKey, err = x509.ParsePKCS1PrivateKey(spkiBlock.Bytes)
		if err != nil {
			log.Print(err.Error())
			return
		}

	}
	publicKey := privateKey.PublicKey
	pubKey := crypt.KeyAsPEMStr(&publicKey)
	log.Printf("This is public key \n%v", pubKey)
	r := gin.Default()
	r.POST("/ping", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Print(err)
			return
		}
		defer c.Request.Body.Close()
		res, err := crypt.RsaOaepDecrypt(string(body), *privateKey)
		var y map[string]interface{}
		err = json.Unmarshal([]byte(res), &y)
		if err != nil {
			log.Print(err)
			return
		}
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
