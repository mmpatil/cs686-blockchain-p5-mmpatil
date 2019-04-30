package data

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

type PeerList struct {
	selfId    int32            `json:"selfId"`
	peerMap   map[string]int32 `json:"peerMap"`
	maxLength int32            `json:"maxLength"`
	mux       sync.Mutex       `json:"mux"`
}

//NewPeerList method will Initialise new PeerList
func NewPeerList(id int32, maxLength int32) PeerList {
	peerList := PeerList{
		selfId:    id,
		peerMap:   make(map[string]int32),
		maxLength: maxLength,
	}
	return peerList
}

//Adds a peer in the PeerList PeerMap
func (peers *PeerList) Add(addr string, id int32) {
	peers.mux.Lock()
	peers.peerMap[addr] = id
	peers.mux.Unlock()
}

//Delete method Deletes a peer from the PeerMap in the PeerList
func (peers *PeerList) Delete(addr string) {
	peers.mux.Lock()
	delete(peers.peerMap, addr)
	peers.mux.Unlock()
}

//Rebalance method Rebalances the PeerMap in the PeerList
func (peers *PeerList) Rebalance() {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	if len(peers.peerMap) > int(peers.maxLength) {
		idList := peers.Sort()
		idList2 := peers.calculateRebalance(idList, peers.selfId)
		peers.RebalancePeerMap(idList2)
	}
}

//RebalancePeerMap method rebalances the PeerMap
func (peers *PeerList) RebalancePeerMap(idList []int32) {
	rebalancedPeerMap := make(map[string]int32)
	for k, v := range peers.peerMap {
		if contains(idList, v) {
			rebalancedPeerMap[k] = v
		}
	}
	peers.peerMap = rebalancedPeerMap
}

func contains(s []int32, e int32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//calculateRebalance method is used to calculate the number of nodes before this node
func (peers *PeerList) calculateRebalance(idList []int32, selfID int32) []int32 {
	index := 0
	idListNew := make([]int32, peers.maxLength)
	a := 0
	// create a ring and populate it with some values
	ring1 := ring.New(len(idList))
	for i := 0; i < ring1.Len(); i++ {
		ring1.Value = idList[i]
		if idList[i] == selfID {
			index = i
		}
		ring1 = ring1.Next()
	}

	ring2 := ring1
	ring1 = ring1.Move(index + 1)
	k := 1
	ring1.Do(func(x interface{}) {
		if k <= int(peers.maxLength/2) {
			idListNew[a] = x.(int32)
			a++
		}
		k++
	})
	f := int(peers.maxLength / 2)
	ring2 = ring2.Move(index)

	for j := 0; j < f; j++ {
		ring2 = ring2.Prev()
		idListNew[a] = ring2.Value.(int32)
		a++
	}
	return idListNew
}

//Sort method Sorts the PeerMap in the PeerList
func (peers *PeerList) Sort() []int32 {
	peerMap := peers.peerMap
	var idList []int32
	for _, value := range peerMap {
		idList = append(idList, int32(value))
	}
	idList = append(idList, peers.selfId)
	sort.Slice(idList, func(i, j int) bool { return idList[i] < idList[j] })
	return idList
}

//Show method shows all addresses and their corresponding IDs. For example, it returns "This is PeerMap: \n addr=127.0.0.1, id=1".
func (peers *PeerList) Show() string {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	result := "This is PeerMap:\n"
	for key, value := range peers.peerMap {
		k := "addr=" + key
		v := ", id=" + strconv.Itoa(int(value))
		result = result + k
		result = result + v
		result = result + "\n"
	}
	fmt.Println(result)
	return result
}

//
func (peers *PeerList) Register(id int32) {
	// peers.mux.Lock()
	// defer peers.mux.Unlock()
	peers.selfId = id
}

//Copy method copies all contents of the PeerMap in the PeerList
func (peers *PeerList) Copy() map[string]int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	newPeerMap := make(map[string]int32)
	for key, value := range peers.peerMap {
		newPeerMap[key] = value
	}
	return newPeerMap
}

//
func (peers *PeerList) GetSelfId() int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	return peers.selfId
}

//PeerMapToJson method Convert PeerMap in PeerList to Json
func (peers *PeerList) PeerMapToJson() (string, error) {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	if len(peers.peerMap) > 0 {
		jsonString, err := json.Marshal(peers.peerMap)
		if err != nil {
			return "", err
		} else {
			return string(jsonString), nil
		}
	}
	return "", nil
}

//InjectPeerMapJson function will take a JSON of a peerMap, and you insert every entry of that peerMap into your own PeerList.
func (peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, selfAddr string) {
	//fmt.Println("in InjectPeerMapJson : selfAddr : ", selfAddr)
	peerMapNew := make(map[string]int32)
	err := json.Unmarshal([]byte(peerMapJsonStr), &peerMapNew)
	if err == nil {
		peers.mux.Lock()
		for key, value := range peerMapNew {
			if key != selfAddr {
				//fmt.Println("in InjectPeerMapJson : value : ", value)
				peers.peerMap[key] = value
			}
		}
		peers.mux.Unlock()
	}
}
