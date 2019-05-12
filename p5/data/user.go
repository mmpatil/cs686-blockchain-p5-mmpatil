package data

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type User struct {
	nationalId  string          `json:"nationalId"`
	privateKey  *rsa.PrivateKey `json:"privateKey"`
	publicKey   *rsa.PublicKey  `json:"publicKey"`
	candidateId int             `json:"candidateId"`
}

type RequestResponse struct {
	user      User           `json:"user"`
	signature []byte         `json:"signature"`
	publicKey *rsa.PublicKey `json:"publicKey"`
}

func RegisterUser(authenticationServerRegister string) []byte {
	fmt.Println("In RegisterUser")
	response, err := http.Get(authenticationServerRegister)
	if err != nil {
		log.Fatal("Error in RegisterUser:", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	//take the public private key and assign to the user
	return body
}

func (user *User) Vote(candidateId int, authenticationServerAddress string, publicKeyServer *rsa.PublicKey) {
	user.candidateId = candidateId
	dummyUser := User{nil, nil, user.publicKey, candidateId}
	userString, err := json.Marshal(dummyUser)
	if err != nil {
		log.Fatal("Error in vote:", err)
	}
	sig, err := GenerateSignature(userString, user.privateKey)
	if err != nil {
		log.Fatal("Error in Vote while Generating signature")
	}
	request := RequestResponse{*user, sig, user.publicKey}
	requestByteArray, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error in Vote in converting RequestResponse object to json ")
	}
	resp, err := http.Post(authenticationServerAddress, "application/json; charset=UTF-8", strings.NewReader(string(requestByteArray)))
	if err != nil {
		log.Fatal("Error in post request to vote")
	}

	fmt.Println("Response of the Post request in Vote:", resp)
}

func (user *User) GetVoteDetails(peerAddress string) {
	dummyUser := User{nil, nil, user.publicKey, nil}
	userString, err := json.Marshal(dummyUser)
	if err != nil {
		log.Fatal("Error in GetVoteDetails: converting to json", err)
	}
	sig, err := GenerateSignature(userString, user.privateKey)
	if err != nil {
		log.Fatal("Error in GetVoteDetails: in GenerateSignature", err)
	}
	request := RequestResponse{*user, sig, user.publicKey}
	requestByteArray, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error in GetVoteDetails : ", err)
	}

	resp, err := http.Post(peerAddress, "application/json; charset=UTF-8", strings.NewReader(string(requestByteArray)))
	if err != nil {
		log.Fatal("Error in post request to vote")
	}
	fmt.Println("Response of post request:", resp)
}
