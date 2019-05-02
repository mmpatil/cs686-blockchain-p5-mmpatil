package p5

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"golang.org/x/crypto/sha3"
	"log"
)

type KeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

//generate public-private key pair
func GenerateKeyPair() *rsa.PrivateKey {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		log.Fatal("Error in generating key-value pair, error is", err)
	}
	return privateKey
}

func (keyPair *KeyPair) GenerateSignature(message []byte) ([]byte, error) {
	hashed := sha3.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, keyPair.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		log.Fatal("Error in generating signature", err)
		return nil, err
	}
	return signature, nil
}

func VerifySignature() {
	rsa.VerifyPKCS1v15()
}
