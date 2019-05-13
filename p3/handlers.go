package p3

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"

	"../blockpackage"
	"../p1"
	"../p2"
	"./data"
	"github.com/gorilla/mux"
)

var TA_SERVER = "http://localhost:6688"
var AUTH_SERVER_REGISTER = ""
var REGISTER_SERVER = TA_SERVER + "/peer"
var BC_DOWNLOAD_SERVER = TA_SERVER + "/upload"
var FIRST_ADDR = "http://localhost:6686"
var REUSE_ADDR = "http://localhost:"
var SELF_ADDR string
var MakeParent bool
var SBC data.SyncBlockChain
var Peers data.PeerList
var ReceivingBlockHeight int32
var ifStarted bool
var port int32
var difficulty int

//
func init() {
	if os.Args[2] == "peer" {
		body := os.Args[1]
		SELF_ADDR = REUSE_ADDR + body
		SBC = data.NewBlockChain()
		difficulty = 2
		if SELF_ADDR == FIRST_ADDR {
			fmt.Println("Port Number:", os.Args[1])
			mpt := getMPT()
			findingNonce := false
			for findingNonce == false {
				nonce := makeNonce(difficulty)
				str := "genesis" + nonce + mpt.Root
				if isProofOfWork(str, difficulty) {
					fmt.Println("Nonce found...")
					b1 := SBC.Initial(mpt, nonce)
					SBC.Insert(b1)
					findingNonce = true
				}
			}
		}
	}
	if os.Args[2] == "auth" {

	}
	if os.Args[2] == "client" {

	}
}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start")
	if ifStarted == false {
		ifStarted = true
		id := Register()
		Peers = data.NewPeerList(id, 32)
		fmt.Fprintf(w, "%s\n", strconv.Itoa(int(id)))
		if SELF_ADDR != FIRST_ADDR {
			Download()
			Peers.Add(FIRST_ADDR, int32(6686))
		}
		go StartHeartBeat()
		StartTryingNonces()
	}
}

func StartClient(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", "Client started....")

}

func StartAuthServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", "Authentication server started....")

}

// Display peerList and sbc
func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

// Register to TA's server, get an ID
func Register() int32 {
	fmt.Println("Register")
	body := os.Args[1]
	port1, _ := strconv.Atoi(body)
	port = int32(port1)
	id, err := strconv.Atoi(string(body))
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return int32(id)
}

// Download blockchain from TA server
func Download() {
	fmt.Println("Download")
	resp, err := http.Get("http://localhost:6686/upload")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) //blockChainJson
	if err != nil {
		log.Fatal(err)
	}
	SBC.UpdateEntireBlockChain(string(body))
}

// Upload blockchain to whoever called this method, return jsonStr
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		data.PrintError(err, "Upload")
	}
	fmt.Fprint(w, blockChainJson)
	//UploadFirstBlock(w, r)
}

// Upload blockchain to whoever called this method, return jsonStr
func UploadFirstBlock(w http.ResponseWriter, r *http.Request) {
	nbc := data.NewBlockChain()
	gbl, _ := SBC.Get(1)
	nbc.Insert(gbl[0])
	blockChainJson, err := nbc.BlockChainToJson()
	if err != nil {
		log.Println("in Err of Upload Genesis")
	}
	_, err = fmt.Fprint(w, blockChainJson)
	if err != nil {
		log.Println("in Err of Upload Genesis writing response")
	}
}

// Upload a block to whoever called this method, return jsonStr
// /block/height/hash
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	height, err := strconv.Atoi(vars["height"])
	if err != nil {
		returnCode500(w, r)
	} else {
		hash := vars["hash"]
		newBlock, exists := SBC.GetBlock(int32(height), hash)
		if exists == true {
			newBlockJson, err := newBlock.EncodeToJSON()
			if err != nil {
				fmt.Println("\n Returning 500 Server error")
				returnCode500(w, r)
			}
			//w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, newBlockJson)
		}
		if exists == false {
			fmt.Println("\n Returning 204 not exists")
			returnCode204(w, r)
		}
	}
}

func returnCode500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func returnCode204(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Block does not exists", http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
	http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
}

// HeartBeatReceive methhReceived a heartbeat
func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	var newheartBeat data.HeartBeatData
	err := json.Unmarshal(bytes, &newheartBeat)
	if err != nil {
		data.PrintError(err, "HeartBeatReceive")
		return
	}
	ProcessHeartBeat(newheartBeat)
	Forward(newheartBeat)
}

// //
// func MPTReceive(w http.ResponseWriter, r *http.Request) p1.MerklePatriciaTrie {
// 	//bytes, _ := ioutil.ReadAll(r.Body)
// 	vars := mux.Vars(r)
// 	insertString, err := vars["message"]
// 	key, _ := vars["key"]
// 	mpt := CreateMPT(key, insertString)
// 	return mpt
// }

// //
// func CreateMPT(key string, message string) p1.MerklePatriciaTrie {
// 	mpt := p1.MerklePatriciaTrie{}
// 	mpt.Initial()
// 	mpt.Insert(key, message)
// 	return mpt
// }

//
func StartTryingNonces() {
	foundNonce := false
	for {
		parentBlocks, exists := SBC.GetLatestBlocks()
		if exists == true {
			foundNonce = false
			parentBlock := parentBlocks[0]
			parentHash := parentBlock.Header.Hash
			height := parentBlock.Header.Height
			mpt := getMPT()
			str := ""
			str = str + parentHash
			for foundNonce == false {
				nonce := makeNonce(difficulty)
				str = str + nonce
				str = str + mpt.Root
				if isProofOfWork(str, difficulty) {
					b1 := SBC.GenBlock(mpt, nonce, height)
					blockjson, _ := b1.EncodeToJSON()
					if ReceivingBlockHeight != SBC.GetLatestHeight() {
						SBC.Insert(b1)
						HeartBeat(blockjson, true)
					}
					foundNonce = true
				}
			}
		}
	}
}

//
func isProofOfWork(str string, difficulty int) bool {
	bytes := sha3.Sum256([]byte(str))
	generatedsha := hex.EncodeToString(bytes[:])
	toCalculate := string(generatedsha[:difficulty])
	compareTo := ""
	for i := 0; i < difficulty; i++ {
		compareTo = compareTo + "0"
	}
	if strings.Compare(toCalculate, compareTo) == 0 {
		return true
	}
	return false
}

//
func makeNonce(lenght int) string {
	bytes := make([]byte, lenght)
	source := rand.NewSource(time.Now().UnixNano())
	for i := range bytes {
		bytes[i] = byte(source.Int63())
	}
	nonce := hex.EncodeToString(bytes)
	return nonce
}

//ProcessHeartBeat method
func ProcessHeartBeat(newheartBeat data.HeartBeatData) {
	var insert bool
	peerMapJson := newheartBeat.PeerMapJson
	registerData := data.NewRegisterData(port, peerMapJson)
	RegisterPeerMap(registerData, newheartBeat)
	if newheartBeat.IfNewBlock == true {
		blockJSON := newheartBeat.BlockJson
		newBlock, _ := blockpackage.DecodeFromJSON(blockJSON)
		ReceivingBlockHeight = newBlock.Header.Height
		parentHash := newBlock.Header.ParentHash
		nonce := newBlock.Header.Nonce
		str := parentHash + nonce + newBlock.Value.Root
		if isProofOfWork(str, difficulty) {
			blockChain := SBC.GetBlockChain()
			height := blockChain.Length
			blocks, _ := blockChain.Get(height) //get current block at a height from blockchain
			for _, block := range blocks {
				if block.Header.Hash == parentHash {
					SBC.Insert(newBlock)
					insert = true
				}
			}
			if insert == false {
				blocks := []blockpackage.Block{}
				i := newBlock.Header.Height - 1
				latestBlocks, _ := SBC.GetLatestBlocks()
				completed := false
				for i > 0 {
					block, exists := AskForBlock(i, newBlock.Header.ParentHash)
					if exists {
						blocks[i] = block
						i--
						for _, latestBlock := range latestBlocks {
							if block.Header.ParentHash == latestBlock.Header.Hash {
								completed = true
								break
							}
						}
					}
					if completed == true {
						break
					}
				}
				if len(blocks) > 0 {
					for _, block := range blocks {
						SBC.Insert(block)
					}
					fmt.Println("Block exists and recovered from another peer")
				} else {
					fmt.Println("Block does not exists")
				}
			}
		}
	}
}

//Forward Method
func Forward(heartBeat data.HeartBeatData) {
	hops := heartBeat.Hops
	if hops > 0 {
		hops = hops - 1
		heartBeat.Hops = hops
		ForwardHeartBeat(heartBeat)
	}
}

//RegisterPeerMap method
func RegisterPeerMap(registerData data.RegisterData, heartBeat data.HeartBeatData) {
	newPeerMapJson := registerData.PeerMapJson
	Peers.InjectPeerMapJson(newPeerMapJson, SELF_ADDR)
	Peers.Add(heartBeat.Addr, heartBeat.Id)
}

// AskForBlock method Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string) (blockpackage.Block, bool) {
	var newBlock blockpackage.Block
	var peersToRemove []string
	newBlock = blockpackage.Block{}
	thePeerMap := Peers.Copy()
	blockExists := false
	if len(thePeerMap) > 0 {
		for k, _ := range thePeerMap {
			resp, err1 := http.Get(k + "block/" + strconv.Itoa(int(height)) + "/" + hash)
			if err1 != nil {
				//fmt.Println("Error")
				continue
			} else {
				if resp.StatusCode == 404 {
					fmt.Println("404", k)
					peersToRemove = append(peersToRemove, k)
					continue
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("here error in reading body")
					continue
				}
				if len(body) > 0 {
					newBlock, _ = blockpackage.DecodeFromJSON(string(body))
					blockExists = true
					break
				}
			}
		}
		if len(peersToRemove) > 0 {
			for _, peer := range peersToRemove {
				Peers.Delete(peer)
			}
		}
	}
	return newBlock, blockExists
}

//
func ForwardHeartBeat(heartBeatData data.HeartBeatData) {
	Peers.Rebalance()
	heartBeatData.Addr = SELF_ADDR
	heartBeatData.Id = Peers.GetSelfId()
	heartBeatData.PeerMapJson, _ = Peers.PeerMapToJson()
	thePeerMap := Peers.Copy()
	heartBeat, _ := json.Marshal(heartBeatData)
	if len(thePeerMap) > 0 {
		for k, _ := range thePeerMap {
			http.Post(k+"/heartbeat/receive", "application/json; charset=UTF-8", strings.NewReader(string(heartBeat)))
		}
	}
}

//HeartBeat method
func StartHeartBeat() {
	for true {
		HeartBeat("", false)
		time.Sleep(5 * time.Second)
	}
}

//StartHeartBeat method
func HeartBeat(blockJSON string, blockGenerate bool) {
	var peersToRemove []string
	jsonString, _ := Peers.PeerMapToJson()
	Peers.Rebalance()
	thePeerMap := Peers.Copy()
	heartBeat := data.PrepareHeartBeatData(&SBC, Peers.GetSelfId(), jsonString, SELF_ADDR, blockGenerate, blockJSON)
	heartBeatData, _ := json.Marshal(heartBeat)
	if len(thePeerMap) > 0 {
		for k, _ := range thePeerMap {
			_, err := http.Post(k+"/heartbeat/receive", "application/json; charset=UTF-8", strings.NewReader(string(heartBeatData)))
			if err != nil {
				peersToRemove = append(peersToRemove, k)
			}
		}
		if len(peersToRemove) > 0 {
			for _, peer := range peersToRemove {
				Peers.Delete(peer)
			}
			Peers.Rebalance()
		}
	}
}

//Canonical func -  Display canonical chain
func Canonical(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In Canonical...")
	canonicalChains := GetCanonicalChains(&SBC)
	fmt.Fprint(w, "Canonical Chain(s) : \n")
	for i, chain := range canonicalChains {
		fmt.Fprint(w, "\nChain #"+strconv.Itoa(i+1))
		_, err := fmt.Fprint(w, "\n", chain.ShowCanonical())
		if err != nil {
			fmt.Fprint(w, "Error in Canonical")
		}
	}
	//
}

// func GetCanonicalChains(SBC *data.SyncBlockChain) []p2.BlockChain {
// 	fmt.Println("In Get canonical chain...")
// 	blocks, _ := SBC.GetLatestBlocks()
// 	fmt.Println("Latest Block:", blocks)
// 	block := blocks[0]
// 	maxHeight := block.Header.Height
// 	canonicalChains := make([]p2.BlockChain, maxHeight)
// 	for i := range blocks {
// 		bc := p2.BlockChain{}
// 		bc.Initial()
// 		canonicalChains[i] = bc
// 	}
// 	for block := range blocks {
// 		canonicalChains[block].Insert(blocks[block])
// 	}
// 	for _, chain := range canonicalChains {
// 		for height := maxHeight - 1; height > 0; height-- {
// 			blockInCanonicalChain, _ := chain.Get(height + 1)
// 			parentBlocksInBlockchain, _ := SBC.Get(height)
// 			for _, parentBlock := range parentBlocksInBlockchain {
// 				if blockpackage.Block(blockInCanonicalChain[0]).Header.ParentHash == blockpackage.Block(parentBlock).Header.Hash {
// 					chain.Insert(parentBlock)
// 				}
// 			}
// 		}
// 	}
// 	return canonicalChains
// }

func GetCanonicalChains(SBC *data.SyncBlockChain) []p2.BlockChain {
	maxHeight := int32(SBC.GetLength())
	blocksAtMaxHeight, _ := SBC.Get(maxHeight)

	canonicalChains := make([]p2.BlockChain, len(blocksAtMaxHeight))
	for i := range blocksAtMaxHeight {
		bc := p2.BlockChain{}
		bc.Initial()
		canonicalChains[i] = bc
	}
	for lastblocks := range blocksAtMaxHeight {
		canonicalChains[lastblocks].Insert(blocksAtMaxHeight[lastblocks])

	}
	for _, chain := range canonicalChains {
		for height := maxHeight - 1; height > 0; height-- {
			existingChildBlocks, _ := chain.Get(int32(height + 1))
			potentialParentBlocks, _ := SBC.Get(int32(height))
			for _, potentialParentBlock := range potentialParentBlocks {
				if blockpackage.Block(existingChildBlocks[0]).Header.ParentHash == blockpackage.Block(potentialParentBlock).Header.Hash {
					chain.Insert(potentialParentBlock)
				}
			}

		}
	}
	return canonicalChains
}

//
func getMPT() p1.MerklePatriciaTrie {
	mpt := p1.MerklePatriciaTrie{}
	mpt.Initial()
	mpt.Insert("do", "verb")
	mpt.Insert("dog", "puppy")
	mpt.Insert("doge", "coin")
	mpt.Insert("horse", "stallion")
	return mpt
}
