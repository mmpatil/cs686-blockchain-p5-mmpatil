package p2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	block1 "../blockpackage"
	"../p1"
	"golang.org/x/crypto/sha3"
)

//Blockchain structure
type BlockChain struct {
	Chain  map[int32][]block1.Block
	Length int32
}

//Get method to get block at a particular height in blockchain
func (blockchain *BlockChain) Get(height int32) ([]block1.Block, bool) {
	if height <= blockchain.Length && height > 0 {
		//_, ok := blockchain.Chain[height]
		//fmt.Println("OK ? :", ok /*, " : ", val*/)
		return blockchain.Chain[height], true
	}
	return nil, false
}

// Insert func takes a block, use the height to insert blockhash , but ignore if hash alrady present
func (blockchain *BlockChain) Insert(block block1.Block) {

	blockHeight := block.Header.Height
	isValidBlock := false

	if blockchain.Length == 0 && blockHeight == 1 && block.Header.ParentHash == "genesis" {
		isValidBlock = genesis(blockchain, block, blockHeight)
	} else if blockHeight > 0 && blockHeight <= blockchain.Length { //adding fork
		isValidBlock = addFork(blockchain, block, blockHeight)
	} else if blockHeight > 0 && blockHeight >= blockchain.Length+1 { //can be any height greater than chain length
		isValidBlock = addLength(blockchain, block, blockHeight)
	}
	if isValidBlock == true {
		fmt.Println("block added to blockchain")
	} else {
		fmt.Println("invalid block, will not add to blockchain !!!!")
	}

}

//genesis func creates the 1st block of blockchain
func genesis(blockchain *BlockChain, block block1.Block, blockHeight int32) bool {
	blockchain.Chain[blockHeight] = append(blockchain.Chain[blockHeight], block)
	blockchain.Length++
	return true
}

//addLength function adds a block such that it increases the length of block
func addLength(blockchain *BlockChain, block block1.Block, blockHeight int32) bool {

	blockchain.Chain[blockHeight] = append(blockchain.Chain[blockHeight], block)
	blockchain.Length = blockHeight
	return true

}

//addFork method adds a block at previously known height
func addFork(blockchain *BlockChain, block block1.Block, blockHeight int32) bool {
	blockList := blockchain.Chain[blockHeight]

	isBlockCorrect := true
	for i := range blockList {
		if blockList[i].Header.Hash == block.Header.Hash {
			isBlockCorrect = false
			break
		}
	}
	if isBlockCorrect == false {
		return false
	}

	blockList = append(blockList, block)
	blockchain.Chain[blockHeight] = blockList // replacing with new blocklist
	return true

}

//EncodeToJSON method encode the blockchain in json format
func (blockchain *BlockChain) EncodeToJSON() (string, error) {
	jsonBlockChainString := "["
	for _, v := range blockchain.Chain {
		for _, b := range v {
			block, err := b.EncodeToJSON()
			if err == nil {
				jsonBlockChainString += block
				jsonBlockChainString += ","
			}
		}
	}
	jsonBlockChainString = jsonBlockChainString[:len(jsonBlockChainString)-1]
	jsonBlockChainString += "]"
	return jsonBlockChainString, nil
}

//Initial methos initializes a blockchain
func (blockchain *BlockChain) Initial() {
	blockchain.Chain = make(map[int32][]block1.Block)
	blockchain.Length = 0
}

//DecodeFromJSON method of blockchain decodes a json to form a blockchain
func DecodeFromJSON(blockchain *BlockChain, jsonString string) {
	//	blockchain := new(BlockChain)
	//	blockchain.Initial()
	var blocks []block1.BlockJson
	err := json.Unmarshal([]byte(jsonString), &blocks)
	if err == nil {
		for _, block := range blocks {
			jsonBlock, err := json.Marshal(block)
			if err == nil {
				block, _ := block1.DecodeFromJSON(string(jsonBlock))
				blockchain.Insert(block)
			}
		}
	}
}

//GetBlock method Gets a block at a height and a particular hash from Blockchain
func (blockchain *BlockChain) GetBlock(height int32, hash string) (block1.Block, bool) {
	blockEmpty := block1.Block{}
	if height <= blockchain.Length {
		blocks, _ := blockchain.Get(height)
		for _, block := range blocks {
			if block.Header.Hash == hash {
				return block, true
			}
		}
	}
	return blockEmpty, false
}

//CheckParentHash method
func (blockchain *BlockChain) CheckParentHash(insertBlock block1.Block) bool {
	blocks, _ := blockchain.Get(blockchain.Length)
	for _, block1 := range blocks {
		if block1.Header.Hash == insertBlock.Header.Hash {
			return true
		}
	}
	return false
}

//todotodo
func (blockchain *BlockChain) InitialBlock(mpt p1.MerklePatriciaTrie, nonce string) block1.Block {
	block := block1.Block{}
	if blockchain.Length == 0 {
		block.Initial(1, "genesis", mpt, nonce)
	}
	return block
}

//
func NewBlockChain() BlockChain {
	blockChain := BlockChain{}
	blockChain.Initial()
	return blockChain
}

//
func (blockchain *BlockChain) Show() string {
	rs := ""
	var idList []int
	for id := range blockchain.Chain {
		idList = append(idList, int(id))
	}
	sort.Ints(idList)
	for _, id := range idList {
		var hashs []string
		for _, block := range blockchain.Chain[int32(id)] {
			hashs = append(hashs, block.Header.Hash+"<="+block.Header.ParentHash)
		}
		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	sum := sha3.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is the BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs
	return rs
}

//GetLatestBlocks method gets the blocks at the maximum height of the blockchain
func (blockchain *BlockChain) GetLatestBlocks() ([]block1.Block, bool) {
	if blockchain.Length > 0 {
		return blockchain.Get(blockchain.Length)
	}
	return nil, false
}

//GetParentBlock method
func (blockchain *BlockChain) GetParentBlock(block block1.Block) (block1.Block, bool) {
	block1 := block1.Block{}
	height := block.Header.Height
	parentHash := block.Header.ParentHash

	blocks, exists := blockchain.Get(height - 1)
	if exists == true {
		for _, block := range blocks {
			if reflect.DeepEqual(block.Header.Hash, parentHash) {
				return block, true
			}
		}
	}
	return block1, false
}

//
func (blockchain *BlockChain) ShowCanonical() string {
	rs := ""
	var idList []int
	for id := range blockchain.Chain {
		idList = append(idList, int(id))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(idList)))
	for _, id := range idList {
		var hashs []string
		for _, block := range blockchain.Chain[int32(id)] {
			hashs = append(hashs, "height="+strconv.Itoa(int(block.Header.Height))+
				", timestamp="+strconv.Itoa(int(block.Header.Timestamp))+
				", hash="+block.Header.Hash+
				", parentHash="+block.Header.ParentHash+
				", size="+strconv.Itoa(int(block.Header.Size)))
		}
		sort.Strings(hashs)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}

	return rs
}
