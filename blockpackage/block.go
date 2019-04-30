package blockpackage

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"../p1"
	"golang.org/x/crypto/sha3"
)

//Header of a block
type Header struct {
	Height     int32
	Timestamp  int64
	Hash       string
	ParentHash string
	Nonce      string
	Size       int32
}

//Block of blockchain
type Block struct {
	Header Header
	Value  p1.MerklePatriciaTrie `json:"value"`
}

//The BlockJson to create
type BlockJson struct {
	Height     int32             `json:"height"`
	Timestamp  int64             `json:"timeStamp"`
	Hash       string            `json:"hash"`
	ParentHash string            `json:"parentHash"`
	Size       int32             `json:"size"`
	Nonce      string            `json:"nonce"`
	MPT        map[string]string `json:"mpt"`
}

//Initial block type created.
func (block *Block) Initial(height int32, parentHash string, mpt p1.MerklePatriciaTrie, nonce string) {
	block.Header.Height = height
	block.Header.ParentHash = parentHash
	block.Header.Nonce = nonce
	block.Value = mpt
	block.Header.Timestamp = time.Now().Unix()
	block.Header.Size = int32(len([]byte(block.Value.String())))
	block.Header.Hash = calculateBlockHash(*block)
}

//calculateBlockHash method calculates the block hash
func calculateBlockHash(block Block) string {
	str := string(block.Header.Height) + string(block.Header.Timestamp) + block.Header.ParentHash + block.Value.Root + string(block.Header.Size)
	sum := sha3.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}

//DecodeFromJSON method to create a block from BlockJSon
func DecodeFromJSON(jsonString string) (Block, error) {
	var blockJSON BlockJson
	var block Block
	byteArray := []byte(jsonString)
	err := json.Unmarshal(byteArray, &blockJSON)
	if err == nil {
		block.Header.Height = blockJSON.Height
		block.Header.ParentHash = blockJSON.ParentHash
		block.Header.Nonce = blockJSON.Nonce
		block.Header.Timestamp = blockJSON.Timestamp
		block.Header.Size = blockJSON.Size
		block.Header.Hash = blockJSON.Hash
		mpt := new(p1.MerklePatriciaTrie)
		mpt.Initial()
		for k, v := range blockJSON.MPT {
			mpt.Insert(k, v)
		}
		block.Value = *mpt
		return block, nil
	}
	return block, err
}

//EncodeToJSON method
func (block *Block) EncodeToJSON() (string, error) {
	var blockJSON BlockJson
	blockJSON.Hash = block.Header.Hash
	blockJSON.Height = block.Header.Height
	blockJSON.ParentHash = block.Header.ParentHash
	blockJSON.Nonce = block.Header.Nonce
	blockJSON.Timestamp = block.Header.Timestamp
	blockJSON.Size = block.Header.Size
	blockJSON.MPT = block.Value.GetAll()
	byteArray, err := json.Marshal(blockJSON)
	s := string(byteArray)
	return s, err
}

// func main() {
// 	mpt := p1.MerklePatriciaTrie{}
// 	mpt.Initial()
// 	fmt.Println("inserting do")
// 	mpt.Insert("do", "verb")
// 	fmt.Println("inserting dog")
// 	mpt.Insert("dog", "puppy")
// 	fmt.Println("inserting doge")
// 	mpt.Insert("doge", "coin")
// 	fmt.Println("inserting horse")
// 	mpt.Insert("horse", "stallion")
// }
