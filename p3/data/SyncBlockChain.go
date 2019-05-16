package data

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	block1 "../../blockpackage"
	"../../p1"
	"../../p2"
)

//
type SyncBlockChain struct {
	bc  p2.BlockChain
	mux sync.Mutex
}

//
func NewBlockChain() SyncBlockChain {
	return SyncBlockChain{bc: p2.NewBlockChain()}
}

//Get method will get the array of blocks at a particular height
func (sbc *SyncBlockChain) Get(height int32) ([]block1.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Get(height)
}

//Get method will get the array of blocks at a particular height
func (sbc *SyncBlockChain) GetLatestHeight() int32 {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Length
}

//GetBlock method gets a block at a height with a particular parent hash
func (sbc *SyncBlockChain) GetBlock(height int32, hash string) (block1.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetBlock(height, hash)
}

//Insert method will insert a block into the blockchain
func (sbc *SyncBlockChain) Insert(block block1.Block) {
	fmt.Println("Adding block to blockchain : :-)")
	sbc.mux.Lock()
	sbc.bc.Insert(block)
	sbc.mux.Unlock()
}

//CheckParentHash method Checks if the block to be inserted has the right parent hash
func (sbc *SyncBlockChain) CheckParentHash(insertBlock block1.Block) bool {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	sbc.bc.CheckParentHash(insertBlock)
	return true
}

//UpdateEntireBlockChain method Updates the entire blockchain
func (sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) {
	sbc.mux.Lock()
	p2.DecodeFromJSON(&sbc.bc, blockChainJson)
	sbc.mux.Unlock()
}

//BlockChainToJson method will give the
func (sbc *SyncBlockChain) BlockChainToJson() (string, error) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.EncodeToJSON()
}

//GenBlock method
func (sbc *SyncBlockChain) GenBlock(mpt p1.MerklePatriciaTrie, nonce string, height int32, finalizedVotes block1.FinalizedVotes) block1.Block {
	block := block1.Block{}
	//height := sbc.bc.Length
	blocks, found := sbc.Get(height)
	if found != true {
		fmt.Println("in syncBlockChain - GenBlock : no blocks found")
	}
	parentHash := blocks[0].Header.Hash
	parentHeight := blocks[0].Header.Height
	if parentHash != "" && height >= 0 {
		block.Initial(parentHeight+1, parentHash, mpt, nonce, finalizedVotes)
	}
	return block
}

//GetBlockChain method
func (sbc *SyncBlockChain) GetBlockChain() p2.BlockChain {
	return sbc.bc
}

//GetLatestBlocks method
func (sbc *SyncBlockChain) GetLatestBlocks() ([]block1.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetLatestBlocks()
}

func (sbc *SyncBlockChain) GetLength() int {
	blocks, _ := sbc.GetLatestBlocks()
	height := blocks[0].Header.Height
	return int(height)
}

//GetParentBlock method
func (sbc *SyncBlockChain) GetParentBlock(block block1.Block) (block1.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetParentBlock(block)
}

//Initial method
func (sbc *SyncBlockChain) Initial(mpt p1.MerklePatriciaTrie, nonce string, votes block1.FinalizedVotes) block1.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.InitialBlock(mpt, nonce, votes)
}

//Show method
func (sbc *SyncBlockChain) Show() string {
	return sbc.bc.Show()
}

func (sbc *SyncBlockChain) ShowMPT() string {
	blockChain := sbc.GetBlockChain()
	height := sbc.bc.Length
	var jsonMPT []string
	jsonMPT = make([]string, 0)
	for height > 0 {
		blockArray, exists := blockChain.Get(height)
		if exists {
			var tempBlock block1.Block
			var mpt p1.MerklePatriciaTrie
			var mapMPT map[string]string
			tempBlock = blockArray[0]
			mpt = tempBlock.Value
			mapMPT = mpt.GetAll()

			jsonBytes, err := json.Marshal(mapMPT)
			if err != nil {
				log.Fatal("Error in ShowMPT")
			}
			jsonMPT = append(jsonMPT, string(jsonBytes))
			jsonMPT = append(jsonMPT, "\n")
		}
		height--
	}
	justString := strings.Join(jsonMPT, " ")
	return justString
}
