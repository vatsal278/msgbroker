// Package crypt provides functions for encrypting and decrypting messages using RSA-OAEP encryption, and converting RSA keys between various formats.
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

// RsaOaepEncrypt encrypts a message using RSA-OAEP encryption with the given public key.
// Returns the encrypted message as a base64-encoded string.
func RsaOaepEncrypt(secretMessage string, key rsa.PublicKey) (string, error) {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), err
}

// RsaOaepDecrypt decrypts a message using RSA-OAEP decryption with the given private key.
// Returns the decrypted message.
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

// PubKeyAsPEMStr converts an RSA public key to a PEM-encoded string and then base64-encodes the result.
// Returns the base64-encoded PEM string.
func PubKeyAsPEMStr(pubkey *rsa.PublicKey) string {
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubkey),
		},
	))
	return base64.StdEncoding.EncodeToString([]byte(pubKeyPem))
}

// PrivKeyAsPEMStr converts an RSA private key to a PEM-encoded string and then base64-encodes the result.
// Returns the base64-encoded PEM string.
func PrivKeyAsPEMStr(key *rsa.PrivateKey) string {
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	))
	return base64.StdEncoding.EncodeToString([]byte(pubKeyPem))
}

// PEMStrAsPubKey decodes a base64-encoded PEM string containing an RSA public key
// and returns the corresponding *rsa.PublicKey.
func PEMStrAsPubKey(pubKey string) (*rsa.PublicKey, error) {
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

// PEMStrAsPrivKey decodes a base64-encoded PEM string containing an RSA private key
// and returns the corresponding *rsa.PrivateKey.
func PEMStrAsPrivKey(pubKey string) (*rsa.PrivateKey, error) {
	decodeString, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	spkiBlock, _ := pem.Decode(decodeString)
	if spkiBlock == nil || spkiBlock.Type != "RSA PRIVATE KEY" {
		err := errors.New("failed to decode PEM block containing private key")
		return nil, err
	}
	pubInterface, err := x509.ParsePKCS1PrivateKey(spkiBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return pubInterface, err
}
