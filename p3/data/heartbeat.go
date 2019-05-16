package data

import (
	"fmt"
)

type HeartBeatData struct {
	IfNewBlock  bool   `json:"ifNewBlock"`
	Id          int32  `json:"id"`
	BlockJson   string `json:"blockJson"`
	PeerMapJson string `json:"peerMapJson"`
	Addr        string `json:"addr"`
	Hops        int32  `json:"hops"`
}

//NewHeartBeatData method initializes new heartbeatdata instance
func NewHeartBeatData(ifNewBlock bool, id int32, blockJson string, peerMapJson string, addr string) HeartBeatData {
	var newHeartBeat HeartBeatData
	newHeartBeat.IfNewBlock = ifNewBlock
	newHeartBeat.Id = id
	newHeartBeat.BlockJson = blockJson
	newHeartBeat.PeerMapJson = peerMapJson
	newHeartBeat.Addr = addr
	newHeartBeat.Hops = 2
	return newHeartBeat
}

//todotodo
func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, peerMapJson string, addr string, generateBlock bool, blockJSONString string) HeartBeatData {
	newHeartBeatData := NewHeartBeatData(false, selfId, "", peerMapJson, addr)
	//makeNew := 1 //(time.Now().UnixNano() / 1000000000000) % 4
	if generateBlock == true {
		newHeartBeatData.IfNewBlock = true
		// mpt := getMPT()
		// b1 := sbc.GenBlock(mpt)
		// sbc.Insert(b1)
		//blockJSONString, err := b1.EncodeToJSON()
		// if err != nil {
		// 	newHeartBeatData.BlockJson = "{}"

		// } else {
		newHeartBeatData.BlockJson = blockJSONString
		//}
	}
	return newHeartBeatData
}

//
//func getMPT() p1.MerklePatriciaTrie {
//	mpt := p1.MerklePatriciaTrie{}
//	mpt.Insert("do", "verb")
//	mpt.Insert("dog", "puppy")
//	mpt.Insert("doge", "coin")
//	mpt.Insert("horse", "stallion")
//	return mpt
//}

func PrintError(err error, errorString string) {
	fmt.Println("Error", err)
	fmt.Print("  String:", errorString)

}
