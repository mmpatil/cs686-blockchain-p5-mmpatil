package data

import (
	"fmt"
	"sync"
)

type UsersList struct {
	UserMap map[string]User   `json:"UserMap"`
	PKMap   map[string]string `json:"PKMap"`
	mux     sync.Mutex        `json:"mux"`
}

//
//type PublicKeyList struct{
//	PKMap map[*rsa.PublicKey]string    `json:"PKMap"`
//	mux sync.Mutex					   `json:"mux"`
//}

func NewUserList() UsersList {
	userList := UsersList{
		UserMap: make(map[string]User),
		PKMap:   make(map[string]string),
	}
	return userList
}

func (usersList *UsersList) Verify(nationalId string) (User, bool) {
	//verify the national ID
	usersList.mux.Lock()
	defer usersList.mux.Unlock()
	_, err := usersList.UserMap[nationalId]

	newUser := User{}
	if err == true {
		fmt.Println("User already exists returning error!!!!")
		//	usersList.mux.Unlock()
		return newUser, false
	}
	fmt.Println("User does not exists returning a new user")

	privateKey := GenerateKeyPair()
	priv := privateKey
	pub := &privateKey.PublicKey

	newUser = User{
		NationalId: nationalId,
		PrivateKey: priv,
		PublicKey:  pub,
	}
	usersList.UserMap[nationalId] = newUser
	usersList.PKMap[pub.N.String()] = nationalId
	//usersList.mux.Unlock()
	return newUser, true
}

//todo
func (usersList *UsersList) DisplayUserList() map[string]User {
	//usersList.mux.Lock()
	//defer usersList.mux.Unlock()
	return usersList.UserMap
}

func (usersList *UsersList) CopyUsersMap() map[string]User {
	usersList.mux.Lock()
	defer usersList.mux.Unlock()
	newUsersList := make(map[string]User)
	for key, value := range usersList.UserMap {
		newUsersList[key] = value
	}
	return newUsersList
}

func (usersList *UsersList) CopyPKMap() map[string]string {
	usersList.mux.Lock()
	defer usersList.mux.Unlock()
	newPKMap := make(map[string]string)
	for key, value := range usersList.PKMap {
		newPKMap[key] = value
	}
	return newPKMap
}
