package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func RsaOaepEncrypt(secretMessage string, key rsa.PublicKey) (string, error) {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), err
}

func RsaOaepDecrypt(cipherText string, privKey rsa.PrivateKey) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func KeyAsPEMStr(pubkey *rsa.PublicKey) string {
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubkey),
		},
	))
	return base64.StdEncoding.EncodeToString([]byte(pubKeyPem))
}
func PEMStrAsKey(pubKey string) (*rsa.PublicKey, error) {
	decodeString, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	spkiBlock, _ := pem.Decode(decodeString)
	if spkiBlock == nil || spkiBlock.Type != "RSA PUBLIC KEY" {
		err := errors.New("failed to decode PEM block containing public key")
		return nil, err
	}
	pubInterface, err := x509.ParsePKCS1PublicKey(spkiBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return pubInterface, err
}
