package data

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
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

//Generate Signature
func GenerateSignature(message []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha3.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		log.Fatal("Error in generating signature", err)
		return nil, err
	}
	return signature, nil
}

//Verify Signature
func VerifySignature(pub *rsa.PublicKey, sig []byte, message []byte) bool {
	hashed := sha3.Sum256(message)
	err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], sig)

	if err != nil {
		log.Fatal("Error in verify signature")
		return false
	}
	return true
}

func Test() {
	privateKey := GenerateKeyPair()

	priv := privateKey
	pub := &privateKey.PublicKey

	keyPair := KeyPair{priv, pub}

	fmt.Println(keyPair, "\n\n")

	message := "This is manali"

	messageInBytes := []byte(message)

	sig, err := GenerateSignature(messageInBytes, priv)

	if err != nil {
		log.Fatal("Error in signature generation")
	}

	fmt.Println("The signature is:", sig)

	result := VerifySignature(keyPair.publicKey, sig, messageInBytes)

	fmt.Println("Result of comparing:", result)
}
