package data

import "encoding/json"

type RegisterData struct {
	AssignedId  int32  `json:"assignedId"`
	PeerMapJson string `json:"peerMapJson"`
}

//NewRegisterData method will initialize an instance
func NewRegisterData(id int32, peerMapJson string) RegisterData {
	return RegisterData{
		AssignedId:  id,
		PeerMapJson: peerMapJson,
	}
}

//EncodeToJson method will return json string of the object
func (data *RegisterData) EncodeToJson() (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	newStringRegisterData := string(b)
	return newStringRegisterData, nil
}
