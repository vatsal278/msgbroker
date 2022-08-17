package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base32"
	"encoding/base64"
	"encoding/pem"
	"reflect"
	"testing"
)

func Test_RSA_OAEP_Encrypt(t *testing.T) {

	tests := []struct {
		name        string
		requestBody string
		setupFunc   func() (*rsa.PrivateKey, rsa.PublicKey, error)
		validation  func(string, *rsa.PrivateKey, error)
	}{
		{
			name:        "Success::RSA_OAEP_Encrypt",
			requestBody: "Hello World",
			setupFunc: func() (*rsa.PrivateKey, rsa.PublicKey, error) {
				privKey, err := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				return privKey, pubKey, err
			},
			validation: func(x string, key *rsa.PrivateKey, err error) {
				y, _ := RsaOaepDecrypt(x, *key)
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if !reflect.DeepEqual(y, "Hello World") {
					t.Errorf("Want: %v, Got: %v", "Hello World", y)
				}
			},
		},
		{
			name:        "Failure::RSA_OAEP_Encrypt::Public exponent too small",
			requestBody: "Hello World",
			setupFunc: func() (*rsa.PrivateKey, rsa.PublicKey, error) {
				privKey, err := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				pubKey.E = 1
				return privKey, pubKey, err
			},
			validation: func(x string, key *rsa.PrivateKey, err error) {
				if !reflect.DeepEqual(err.Error(), "crypto/rsa: public exponent too small") {
					t.Errorf("Want: %v, Got: %v", "crypto/rsa: public exponent too small", err.Error())
				}
			},
		},
		{
			name:        "Failure::RSA_OAEP_Encrypt::public exponent too large",
			requestBody: "Hello World",
			setupFunc: func() (*rsa.PrivateKey, rsa.PublicKey, error) {
				privKey, err := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				pubKey.E = 1000000000000
				return privKey, pubKey, err
			},
			validation: func(x string, key *rsa.PrivateKey, err error) {
				if !reflect.DeepEqual(err.Error(), "crypto/rsa: public exponent too large") {
					t.Errorf("Want: %v, Got: %v", "crypto/rsa: public exponent too large", err.Error())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKey, pubKey, _ := tt.setupFunc()

			x, err := RsaOaepEncrypt(tt.requestBody, pubKey)
			tt.validation(x, privKey, err)

		})
	}
}

func Test_RSA_OAEP_Decrypt(t *testing.T) {

	tests := []struct {
		name        string
		requestBody string
		setupFunc   func() (*rsa.PrivateKey, string)
		validation  func(string, error)
	}{
		{
			name:        "Success::RSA_OAEP_Decrypt",
			requestBody: "Hello World",
			setupFunc: func() (*rsa.PrivateKey, string) {
				privKey, err := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				x, err := RsaOaepEncrypt("Hello World", pubKey)
				if err != nil {
					t.Log(err.Error())
				}
				return privKey, x
			},
			validation: func(x string, err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", "nil", err.Error())
				}
				if !reflect.DeepEqual(x, "Hello World") {
					t.Errorf("Want: %v, Got: %v", "Hello World", x)
				}
			},
		},
		{
			name:        "Failure:: RSA_OAEP_Decrypt",
			requestBody: "Hello World",
			setupFunc: func() (*rsa.PrivateKey, string) {
				privKey, err := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				x, err := RsaOaepEncrypt("Hello World", pubKey)
				if err != nil {
					t.Log(err.Error())
				}
				privKey.PublicKey.E = 0
				return privKey, x
			},
			validation: func(x string, err error) {
				if !reflect.DeepEqual(err.Error(), "crypto/rsa: public exponent too small") {
					t.Errorf("Want: %v, Got: %v", "crypto/rsa: public exponent too small", err.Error())
				}
			},
		},
		{
			name:        "Failure::RSA_OAEP_Decrypt::Public exponent too large",
			requestBody: "Hello World",
			setupFunc: func() (*rsa.PrivateKey, string) {
				privKey, err := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				x, err := RsaOaepEncrypt("Hello World", pubKey)
				if err != nil {
					t.Log(err.Error())
				}
				privKey.PublicKey.E = 100000000000
				return privKey, x
			},
			validation: func(x string, err error) {
				if !reflect.DeepEqual(err.Error(), "crypto/rsa: public exponent too large") {
					t.Errorf("Want: %v, Got: %v", "crypto/rsa: public exponent too large", err.Error())
				}
			},
		},
		{
			name:        "Failure::RSA_OAEP_Decrypt::Decryption error",
			requestBody: "Hello World",
			setupFunc: func() (*rsa.PrivateKey, string) {
				privKey, err := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				pubKey.E = 100000000000
				x, err := RsaOaepEncrypt("Hello World", pubKey)
				if err != nil {
					t.Log(err.Error())
				}

				return privKey, x
			},
			validation: func(x string, err error) {
				if !reflect.DeepEqual(err.Error(), "crypto/rsa: decryption error") {
					t.Errorf("Want: %v, Got: %v", "crypto/rsa: decryption error", err.Error())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, x := tt.setupFunc()
			y, err := RsaOaepDecrypt(x, *privateKey)

			tt.validation(y, err)

		})
	}
}

func Test_KeyAsPEMStr(t *testing.T) {

	tests := []struct {
		name        string
		requestBody string
		setupFunc   func() (*rsa.PrivateKey, rsa.PublicKey, error)
		validation  func(string, *rsa.PublicKey)
	}{
		{
			name:        "Success::KeyAsPEMStr",
			requestBody: "Hello World",
			setupFunc: func() (*rsa.PrivateKey, rsa.PublicKey, error) {
				privKey, err := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				return privKey, pubKey, err
			},
			validation: func(x string, pubKey *rsa.PublicKey) {
				y, _ := PEMStrAsKey(x)
				if !reflect.DeepEqual(y, pubKey) {
					t.Errorf("Want: %v, Got: %v", pubKey, y)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, pubKey, _ := tt.setupFunc()
			x := KeyAsPEMStr(&pubKey)
			tt.validation(x, &pubKey)

		})
	}
}
func Test_PEMStrAsKey(t *testing.T) {

	tests := []struct {
		name        string
		requestBody string
		setupFunc   func() (string, rsa.PublicKey)
		validation  func(*rsa.PublicKey, *rsa.PublicKey, error)
	}{
		{
			name:        "Success::PEMStrAsKey",
			requestBody: "Hello World",
			setupFunc: func() (string, rsa.PublicKey) {
				privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				x := KeyAsPEMStr(&pubKey)
				return x, pubKey
			},
			validation: func(newPubKey *rsa.PublicKey, pubKey *rsa.PublicKey, err error) {
				if !reflect.DeepEqual(newPubKey, pubKey) {
					t.Errorf("Want: %v, Got: %v", newPubKey, pubKey)
				}
			},
		},
		{
			name:        "Failure::PEMStrAsKey::failed to decode PEM block::Empty Public Key",
			requestBody: "Hello World",
			setupFunc: func() (string, rsa.PublicKey) {
				privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				x := ""
				return x, pubKey
			},
			validation: func(newPubKey *rsa.PublicKey, pubKey *rsa.PublicKey, err error) {
				if !reflect.DeepEqual(err.Error(), "failed to decode PEM block containing public key") {
					t.Errorf("Want: %v, Got: %v", "failed to decode PEM block containing public key", err.Error())
				}
				if reflect.DeepEqual(newPubKey, pubKey) {
					t.Errorf("Want: %v, Got: %v", newPubKey, pubKey)
				}
			},
		},
		{
			name:        "Failure::PEMStrAsKey::failed to decode PEM block::Incorrect PEM type",
			requestBody: "Hello World",
			setupFunc: func() (string, rsa.PublicKey) {
				privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				pubKeyPem := string(pem.EncodeToMemory(
					&pem.Block{
						Type:  "RSA FAKE PUBLIC KEY",
						Bytes: x509.MarshalPKCS1PublicKey(&pubKey),
					},
				))
				return base64.StdEncoding.EncodeToString([]byte(pubKeyPem)), pubKey
			},
			validation: func(newPubKey *rsa.PublicKey, pubKey *rsa.PublicKey, err error) {
				if !reflect.DeepEqual(err.Error(), "failed to decode PEM block containing public key") {
					t.Errorf("Want: %v, Got: %v", "failed to decode PEM block containing public key", err.Error())
				}
				if reflect.DeepEqual(newPubKey, pubKey) {
					t.Errorf("Want: %v, Got: %v", newPubKey, pubKey)
				}
			},
		},
		{
			name:        "Failure::PEMStrAsKey::illegal base64 data",
			requestBody: "Hello World",
			setupFunc: func() (string, rsa.PublicKey) {
				privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
				pubKey := privKey.PublicKey
				pubKeyPem := string(pem.EncodeToMemory(
					&pem.Block{
						Type:  "RSA PUBLIC KEY",
						Bytes: x509.MarshalPKCS1PublicKey(&pubKey),
					},
				))
				return base32.StdEncoding.EncodeToString([]byte(pubKeyPem)), pubKey
			},
			validation: func(newPubKey *rsa.PublicKey, pubKey *rsa.PublicKey, err error) {
				if !reflect.DeepEqual(err.Error(), "illegal base64 data at input byte 684") {
					t.Errorf("Want: %v, Got: %v", "illegal base64 data at input byte 684", err.Error())
				}
				if reflect.DeepEqual(newPubKey, pubKey) {
					t.Errorf("Want: %v, Got: %v", newPubKey, pubKey)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pemStr, pubKey := tt.setupFunc()
			newPubKey, err := PEMStrAsKey(pemStr)
			t.Log(err)
			tt.validation(newPubKey, &pubKey, err)

		})
	}
}
