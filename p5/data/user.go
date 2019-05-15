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
	NationalId  string          `json:"nationalId"`
	PrivateKey  *rsa.PrivateKey `json:"privateKey"`
	PublicKey   *rsa.PublicKey  `json:"publicKey"`
	CandidateId int             `json:"candidateId"`
}

type RequestResponse struct {
	User      User           `json:"user"`
	Signature []byte         `json:"signature"`
	PublicKey *rsa.PublicKey `json:"publicKey"`
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

func (u *User) EncodeToJson() string {
	jsonb, err := json.Marshal(u)

	if err != nil {
		log.Fatal("Error in EncodeToJson of User:", err)
	}
	return string(jsonb)
}

//todo
func (reqRes *RequestResponse) EncodeRequestRespToJson() string {
	jsonBytes, err := json.Marshal(reqRes)

	if err != nil {
		log.Fatal("Error in EncodeToJson RequestResponse:", err)
	}
	return string(jsonBytes)
}

//todo
func (user *User) Vote(candidateId int, authenticationServerAddress string, publicKeyServer *rsa.PublicKey) bool {
	user.CandidateId = candidateId
	dummyUser := User{"", nil, user.PublicKey, candidateId}
	userString, err := json.Marshal(dummyUser)
	if err != nil {
		log.Println("Error in vote:", err)
		return false
	}
	sig, err := GenerateSignature(userString, user.PrivateKey)
	if err != nil {
		log.Println("Error in Vote while Generating signature")
		return false
	}
	request := RequestResponse{*user, sig, user.PublicKey}
	requestByteArray, err := json.Marshal(request)
	if err != nil {
		log.Println("Error in Vote in converting RequestResponse object to json ")
		return false
	}
	resp, err := http.Post(authenticationServerAddress, "application/json; charset=UTF-8", strings.NewReader(string(requestByteArray)))
	if err != nil {
		log.Println("Error in post request to vote")
		return false
	}
	fmt.Println("Response of the Post request in Vote in user.go:", resp)
	return true
}

func (user *User) GetVoteDetails(peerAddress string) {
	dummyUser := User{"", nil, user.PublicKey, 0}
	userString, err := json.Marshal(dummyUser)
	if err != nil {
		log.Fatal("Error in GetVoteDetails: converting to json", err)
	}
	sig, err := GenerateSignature(userString, user.PrivateKey)
	if err != nil {
		log.Fatal("Error in GetVoteDetails: in GenerateSignature", err)
	}
	request := RequestResponse{*user, sig, user.PublicKey}
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
