package p3

import (
	p5 "../p5/data"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var ClientPeersAlive map[string]int32

var userList p5.UsersList
var user p5.User
var SECOND_ADDR = "http://localhost:6687"

//var REGISTRATION_SERVER string

type BodyToSend struct {
	Height int `json:"Height"`
}

//Marshal the body
func MarshalBody(body BodyToSend) string {
	a, _ := json.Marshal(body)
	return string(a)
}

//Starts a client with initial configuration
func StartClient(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "%s\n", "Client started....")
	tpl.ExecuteTemplate(w, "startclient.html", nil)
	ClientPeersAlive = make(map[string]int32)
	ClientPeersAlive[FIRST_ADDR] = int32(6686)
	//ClientPeersAlive.Add(FIRST_ADDR,int32(6686))
	//ClientPeersAlive.Add(SECOND_ADDR,int32(6687))
}

//Starts RegistrationServer/Authentication Server with initial configuration
func StartRegistrationServer(w http.ResponseWriter, r *http.Request) {
	userList = p5.NewUserList()
}

//
func ShowVoteUserC(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "voteinfo.html", nil)
}

//New Client signup
func SignUp(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "signup.html", nil)
}

//Client Signs in
func SignIn(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "signin.html", nil)
}

//Registers a Client with the registration Server/ Authentication Server.
func RegisterClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	nationalId := r.FormValue("nationalId")

	address := REUSE_ADDR + REGISTRATION_SERVER + REGISTER + "/" + nationalId
	fmt.Println("In Register Client. The address used will be:", address)
	resp, err := http.Get(address)
	if err != nil {
		log.Fatal("Error in Get request in RegisterClient:", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusBadRequest {
			fmt.Fprintf(w, "User already exists!!")
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Error in reading the body:", err)
			}

			if len(body) > 0 {
				newUser := p5.User{}
				//fmt.Println("The body is:", string(body))
				err := json.Unmarshal(body, &newUser)
				if err != nil {
					log.Fatal("Error while unmarshal the body:", err)
				}
				user = newUser

				//fmt.Println("The body is:", body)
				fmt.Println("The response is:", newUser)
				jsonString, _ := json.Marshal(user)
				fmt.Fprintf(w, string(jsonString))
			}
		}
	}
}

//NOT USED
func UserRegister(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nationalId := vars["nationalId"]

	fmt.Println("National ID in UserRegister:", nationalId)

	//
	//if err !=nil{
	//	log.Fatal("Error in User Register")
	//}

	//newUser,valid := userList.Verify(string(nationalId))
	newUser, valid := userList.Verify(nationalId)

	if valid == true {
		fmt.Println("New user is:", newUser)
		//jsonstring, err :=json.Marshal(&newUser)
		//if err != nil {
		//	log.Println("Error in marshalling user : err - ", err)
		//}
		jsonString := newUser.EncodeToJson()
		fmt.Println("string(jsonstring) : ", jsonString)
		fmt.Fprint(w, jsonString)
		//fmt.Fprintf(w,string(jsonstring))
	} else {
		http.Error(w, "User not Registered", http.StatusBadRequest)
		fmt.Fprint(w, "Invalid nationalId! NationalId Exists")
		//fmt.Fprintf(w,"Invalid nationalId! NationalId Exists")
	}
}

func DisplayUsers(w http.ResponseWriter, r *http.Request) {

	var usersArray []string
	userMap := userList.CopyUsersMap()

	for k, _ := range userMap {
		usersArray = append(usersArray, k)
	}
	//
	//userMapBytes,err := json.Marshal(userMap)
	//if err!=nil{
	//	log.Fatal("Error in json.Marshal of userMap in DisplayUsers")
	//}

	fmt.Fprintf(w, "%s\n", usersArray)
}

////sends a post request to
//func ShowBlock(w http.ResponseWriter, r *http.Request) {
//	heightString := r.FormValue("height")
//	height, _ := strconv.Atoi(heightString)
//
//	address := REUSE_ADDR + REGISTRATION_SERVER + "/showBlockAtHeight"
//
//	bodyToSend := BodyToSend{
//		Height: height,
//	}
//
//	body := MarshalBody(bodyToSend)
//
//	_, err2 := http.Post(address, "application/json; charset=UTF-8", strings.NewReader(string(body)))
//	if err2 != nil {
//		log.Fatal("Error in Check User response of Post request")
//	}
//}

//When a voter/user signs in its Public-Key and Private-Key pair is verified with the registration server.
func CheckUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	publicKey := r.FormValue("publicKey")
	privateKey := r.FormValue("privateKey")

	//fmt.Println("PublicKey:", publicKey)
	//fmt.Println("PrivateKey:", privateKey)

	jsonString := "{\"nationalId\":\"556655665566\",\"privateKey\":" + privateKey + ",\"publicKey\":" + publicKey + ",\"candidateId\":0}"

	newUser := p5.User{}
	err := json.Unmarshal([]byte(jsonString), &newUser)

	if err != nil {
		log.Fatal("Error in unmarshal CheckUser manually preparing json!!!")
	}

	//jsonString,_ := json.Marshal(newUser)

	address := REUSE_ADDR + REGISTRATION_SERVER + "/check"

	resp, err2 := http.Post(address, "application/json; charset=UTF-8", strings.NewReader(string(jsonString)))

	if err2 != nil {
		log.Fatal("Error in Check User response of Post request")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error in reading the body:", err)
	}
	defer resp.Body.Close()
	if len(body) > 0 {
		respUser := p5.User{}
		fmt.Println("Body is:", string(body))
		err3 := json.Unmarshal(body, &respUser)
		fmt.Println("User in Response:", respUser)
		if err3 != nil {
			log.Fatal("Error in unmarshal error3:", err3)
		}
		user = respUser

		ClientFetchingPeerList()

		type customData struct {
			Title   string
			Members []string
		}

		cd := customData{
			Title:   "Voting Page",
			Members: []string{string(body)},
		}
		tpl.ExecuteTemplate(w, "vote.html", cd)

	} else {
		fmt.Fprintf(w, "%s\n", "Error in reading the body")
	}
}

//Shows the vote details of a user to that user.
func ShowVoteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	publicKey := r.FormValue("publicKey")
	//privateKey := r.FormValue("privateKey")

	//fmt.Println("PublicKey:", publicKey)
	//fmt.Println("PrivateKey:", privateKey)

	jsonString := "{\"nationalId\":\"556655665566\",\"privateKey\":" + "" + ",\"publicKey\":" + publicKey + ",\"candidateId\":0}"

	newUser := p5.User{}
	err := json.Unmarshal([]byte(jsonString), &newUser)

	if err != nil {
		log.Fatal("Error in unmarshal CheckUser manually preparing json!!!")
	}

	//jsonString,_ := json.Marshal(newUser)

	address := REUSE_ADDR + REGISTRATION_SERVER + "/voteInfo"

	_, err2 := http.Post(address, "application/json; charset=UTF-8", strings.NewReader(string(jsonString)))

	if err2 != nil {
		log.Fatal("Error in Check User response of Post request")
	}

}

//User/Voter fetches a peerlist from the peer 6686
func ClientFetchingPeerList() {

	address := FIRST_ADDR + "/getPeerList"

	respPeers, err := http.Get(address)

	if err != nil {
		log.Fatal("Error while getting Peers from FIRST PEER")
	}
	body, err := ioutil.ReadAll(respPeers.Body)
	if err != nil {
		log.Fatal("Error in reading the body:", err)
	}
	defer respPeers.Body.Close()
	fmt.Println("Fetching Peers:", string(body))
	if len(body) > 0 {
		peerListTemp := make(map[string]int32)
		err2 := json.Unmarshal(body, &peerListTemp)
		if err2 != nil {
			log.Fatal("Error in reading the peerList:", err2)
		}
		ClientPeersAlive = peerListTemp
	}
}

//Checks if the new user/voter trying to register already exists
func Check(w http.ResponseWriter, r *http.Request) {

	oldUserList := userList.CopyUsersMap()
	oldPKMap := userList.CopyPKMap()
	newUser := p5.User{}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error in Check - while reading  body - err  ", err)
	}
	defer r.Body.Close()
	err = json.Unmarshal(bytes, &newUser)
	if err != nil {
		log.Fatal("Error in Server Check function unmarshal", err)
	}
	//	fmt.Println("New user:", newUser)
	//	fmt.Println("PKMap:", oldPKMap)

	pub := *(newUser.PublicKey)
	//fmt.Println("pub:",pub.N.String())
	pubString := pub.N.String()
	nationalId, valid := oldPKMap[pubString]
	//fmt.Println("oldPKMap[*(newUser.PublicKey)]:",*(newUser.PublicKey))
	//fmt.Println("National Id is :",nationalId)

	if valid == false {
		//error return
		fmt.Println("The Public Key is not found")
		http.Error(w, "Error user did not exist", http.StatusInternalServerError)
		return
	}

	//check private key
	//	fmt.Println("Checking private key if equal in Check")
	oldUser, exists := oldUserList[nationalId]

	//fmt.Println("User Details:",oldUser)

	if exists == false {
		fmt.Println("User does not exist with this nationalID:", nationalId)
		http.Error(w, "Error user did not exist", http.StatusInternalServerError)
		return
	} else {
		fmt.Println("National ID exists:", oldUser)

		oldPrivateKey := *(oldUser.PrivateKey)

		oldPrivateKeyString := oldPrivateKey.D.String()

		newPrivateKey := *(newUser.PrivateKey)

		newPrivateKeyString := newPrivateKey.D.String()

		if oldPrivateKeyString == newPrivateKeyString {
			fmt.Println("Private Key is equal!!")
			//userBytes,err2 := json.Marshal(oldUser)
			//if err2 != nil{
			//	log.Fatal("This should never reach:",err2)
			//}
			fmt.Println("Valid User!!!!!!!!!!!!!!!!!!! Returning 200 OK")
			//w.WriteHeader(http.StatusOK)
			jsonString := oldUser.EncodeToJson()
			//	fmt.Println("string(jsonstring) : ", jsonString)
			fmt.Fprint(w, jsonString)
			//fmt.Println("sending Response:",userBytes)
			//fmt.Fprint(w,userBytes)
		} else {
			http.Error(w, "Error user did not exist", http.StatusInternalServerError)
			return
		}
	}
}

//Allows a valid user/voter to vote
func ClientVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	candidateIdString := r.FormValue("candidateId")

	candidateId, err := strconv.Atoi(candidateIdString)

	if err != nil {
		log.Fatal("Error in fetching the candidateID")
	}
	user.CandidateId = candidateId

	//userBytes := user.EncodeToJson()

	if len(ClientPeersAlive) > 0 {
		for k, _ := range ClientPeersAlive {

			address := k + "/vote"
			if user.Vote(candidateId, address, user.PublicKey) {
				break
			} else {
				continue
			}
			//_,err :=http.Post(k+"/vote", "application/json; charset=UTF-8", strings.NewReader(string(userBytes)))
			//if err != nil{
			//	continue
			//}
			//break
		}
	}

	_, err2 := fmt.Fprintf(w, "After voting page")
	if err2 != nil {
		log.Fatal("Error in ClientVote FprintF:", err2)
	}
}

//todo
func VoteDetails(w http.ResponseWriter, r *http.Request) {
	_, err2 := fmt.Fprintf(w, "Vote Details Page!!!!")
	if err2 != nil {
		log.Fatal("Error in ClientVote FprintF:", err2)
	}
}
