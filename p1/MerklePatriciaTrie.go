package p1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/crypto/sha3"
)

type Flag_value struct {
	encoded_prefix []uint8
	value          string
}

type Node struct {
	node_type    int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value   Flag_value
}

type MerklePatriciaTrie struct {
	db   map[string]Node
	Root string
}

func toByteArray(s string) []uint8 {
	mySlice := []uint8(s)
	return mySlice
}

// string_to_hex_array converts string into its corresponding hex value
func string_to_hex_array(s string) []uint8 {
	uint8Array := toByteArray(s)
	hexArray := ascii_to_hex(uint8Array)
	return hexArray
}

func compact_encode(hex_array []uint8) []uint8 {
	var term int
	if hex_array[len(hex_array)-1] == 16 {
		term = 1
	}
	if term == 1 {
		hex_array = hex_array[:len(hex_array)-1]
	}
	oddlen := len(hex_array) % 2
	var flags []uint8
	value := uint8(2*term + oddlen)
	flags = append(flags, value)
	if oddlen > 0 {
		hex_array = append(flags, hex_array...)
	} else {
		flags = append(flags, uint8(0))
		hex_array = append(flags, hex_array...)
	}
	//fmt.Println("hex array ", hex_array)
	ascii_array := []uint8{}

	for i := 0; i < len(hex_array); i += 2 {
		ascii_array = append(ascii_array, 16*hex_array[i]+hex_array[i+1])
	}
	return ascii_array
}

func hexToString(hex_array []uint8) string {
	ascii_array := []uint8{}
	for i := 0; i < len(hex_array); i += 2 {
		ascii_array = append(ascii_array, 16*hex_array[i]+hex_array[i+1])
	}
	s := string(ascii_array)
	return s
}

func ascii_to_hex(encoded_arr []uint8) []uint8 {
	hex_array := []uint8{}
	for _, value := range encoded_arr {
		div, mod := value/16, value%16
		hex_array = append(hex_array, div)
		hex_array = append(hex_array, mod)
	}
	return hex_array
}

// If Leaf, ignore 16 at the end
func compact_decode(encoded_arr []uint8) []uint8 {
	hex_array := ascii_to_hex(encoded_arr)
	flag := hex_array[0]
	switch flag {
	case 0:
		hex_array = hex_array[2:]
	case 1:
		hex_array = hex_array[1:]
	case 2:
		hex_array = hex_array[2:]
	case 3:
		hex_array = hex_array[1:]
	}
	return hex_array
}

// check node for its type
func check_node_type(node Node) int {
	switch node.node_type {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	}
	return 0
}

//check if the node is leaf or extension
func check_if_leaf_or_extension(node Node) ([]uint8, uint8) {
	hexPrefixArray := []uint8{}
	var nodeType uint8
	nodeType = 99
	if check_node_type(node) == 2 {
		hexPrefixArray = ascii_to_hex(node.flag_value.encoded_prefix)
		if len(hexPrefixArray) != 0 {
			nodeType = hexPrefixArray[0]
		}
	}
	return hexPrefixArray, nodeType
}

// To get the value of the key from mpt starter to call recursive function
func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	//fmt.Println("VALUE to be searched:", key)
	pathLeft := string_to_hex_array(key)
	currentNode := mpt.db[mpt.Root]
	finalValue, err := mpt.GetRecursive(currentNode, pathLeft)

	if err != nil {
		return "", err
	} else {
		return finalValue, nil
	}
}

// To get the value of the key from mpt
func (mpt *MerklePatriciaTrie) GetRecursive(currentNode Node, pathLeft []uint8) (string, error) {
	if len(pathLeft) < 0 {
		return "", errors.New("path_not_found")
	} else if len(pathLeft) == 0 { // if path left == 0 value has to be in branch or leaf below
		if check_node_type(currentNode) == 1 { // if branch
			if currentNode.branch_value[16] != "" {
				return currentNode.branch_value[16], nil
			}
		} else if check_node_type(currentNode) == 2 {
			_, nodeType := check_if_leaf_or_extension(currentNode)
			if nodeType == 2 || nodeType == 3 { //if leaf
				if len(compact_decode(currentNode.flag_value.encoded_prefix)) == 0 {
					return currentNode.flag_value.value, nil
				}
			}
		}
	} else if len(pathLeft) > 0 { //path left > 0
		if check_node_type(currentNode) == 1 { //branch node make a recursive call
			if currentNode.branch_value[pathLeft[0]] != "" {
				hash := currentNode.branch_value[pathLeft[0]]
				pathLeft = pathLeft[1:]
				return mpt.GetRecursive(mpt.db[hash], pathLeft)
			}
		} else if check_node_type(currentNode) == 2 {
			_, nodeType := check_if_leaf_or_extension(currentNode)
			if nodeType == 0 || nodeType == 1 { // if extension make a recursive call
				nodePath := compact_decode(currentNode.flag_value.encoded_prefix)
				counter := get_common_lenght(nodePath, pathLeft)
				if counter == (len(nodePath)) {
					pathLeft = pathLeft[counter:]
					return mpt.GetRecursive(mpt.db[currentNode.flag_value.value], pathLeft)
				}
			} else if nodeType == 2 || nodeType == 3 { //if leaf path not found
				nodePath := compact_decode(currentNode.flag_value.encoded_prefix)
				if len(pathLeft) < len(nodePath) {
					return "", errors.New("path_not_found")
				}
				counter := get_common_lenght(nodePath, pathLeft)
				if counter == (len(nodePath)) {
					pathLeft = pathLeft[counter:]
					if len(pathLeft) == 0 {
						return currentNode.flag_value.value, nil
					} else {
						return mpt.GetRecursive(mpt.db[currentNode.flag_value.value], pathLeft)
					}
				}
			}
		}
	}
	return "", errors.New("path_not_found")
}

// gets the common lenght between two slice
func get_common_lenght(nodePath []uint8, pathLeft []uint8) int {
	counter := 0
	for i := 0; i < len(nodePath); i++ {
		if nodePath[i] == pathLeft[i] {
			counter++
		} else {
			break
		}
	}
	return counter
}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	//	fmt.Println("KEY being inserted:", key)
	//	fmt.Println("VALUE being inserted:", new_value)
	hex_key_array := string_to_hex_array(key)
	if mpt.Root == "" {
		hex_key_array = append(hex_key_array, 16)
		nodeL := Node{}
		nodeL.node_type = 2
		nodeL.flag_value.encoded_prefix = compact_encode(hex_key_array)
		nodeL.flag_value.value = new_value
		mpt.Root = mpt.InsertRecursive(nodeL, nil, "")
	} else {
		root_node_hash := mpt.Root
		mpt.Root = mpt.InsertRecursive(mpt.db[root_node_hash], hex_key_array, new_value)
	}
}

func (mpt *MerklePatriciaTrie) InsertRecursive(node Node, leftPath []uint8, newValue string) string {
	currNode := node
	old_hash := currNode.hash_node()
	if len(leftPath) == 0 && newValue == "" { //insert in to the db and return hash
		hash := currNode.hash_node()
		if mpt.db == nil {
			mpt.db = make(map[string]Node)
		}
		mpt.db[hash] = currNode
		return hash
	} else if len(leftPath) == 0 { // path left == 0 next node has to be branch or leaf
		if currNode.node_type == 1 { // if branch update or insert values in [16] field
			currNode.branch_value[16] = newValue
			hash := currNode.hash_node()
			delete(mpt.db, old_hash)
			mpt.db[hash] = currNode
			return hash
		} else if currNode.node_type == 2 {
			hexPrefixArray, nodeType := check_if_leaf_or_extension(currNode)
			if nodeType == 2 || nodeType == 3 { //if leaf check if the encoded prefix is same to update value or create branch and leaf
				hexDecodeArray := compact_decode(currNode.flag_value.encoded_prefix)
				if reflect.DeepEqual(hexPrefixArray, []uint8{2, 0}) {
					currNode.flag_value.value = newValue
					hash := currNode.hash_node()
					delete(mpt.db, old_hash)
					mpt.db[hash] = currNode
					return hash
				} else if len(hexDecodeArray) > 0 {
					hexDecodeArray = append(hexDecodeArray, 16)

					nodeL := createLeaf() //create leaf
					nodeL.flag_value.value = currNode.flag_value.value
					nodeL.flag_value.encoded_prefix = compact_encode(hexDecodeArray[1:])

					nodeB := createBranch() //create branch
					nodeB.branch_value[16] = newValue
					nodeB.branch_value[hexDecodeArray[0]] = mpt.InsertRecursive(nodeL, nil, "")
					hash := nodeB.hash_node()
					delete(mpt.db, old_hash)
					mpt.db[hash] = nodeB
					return hash
				}
			}
		}
	} else { // path left > 1
		if check_node_type(currNode) == 2 {
			_, nodeType := check_if_leaf_or_extension(currNode)
			decodedHexArray := compact_decode(currNode.flag_value.encoded_prefix)
			if nodeType == 2 || nodeType == 3 { //if leaf
				if reflect.DeepEqual(decodedHexArray, leftPath) { //same
					currNode.flag_value.value = newValue
					leftPath = nil
					delete(mpt.db, old_hash)
					hash := currNode.hash_node()
					mpt.db[hash] = currNode
					return hash
				} else if check_if_common_path_exists(leftPath, decodedHexArray) { //common path exists
					counter := 0
					for i := 0; i < len(decodedHexArray); i++ {
						if i == len(decodedHexArray) || i == len(leftPath) {
							break
						}
						if decodedHexArray[i] == leftPath[i] {
							counter = counter + 1
						} else {
							break
						}
					}
					commonPath := leftPath[:(counter)]
					leftPath = leftPath[counter:]
					leftNibble := decodedHexArray[counter:]

					nodeE := createExtension() //create extension node and branch
					nodeE.flag_value.encoded_prefix = compact_encode(commonPath)

					nodeB := createBranch() //create branch
					if len(leftNibble) > 0 {
						leftNibble = append(leftNibble, 16)
					}
					if len(leftPath) > 0 {
						leftPath = append(leftPath, 16)
						nodeL1 := createLeaf()
						index := leftPath[0]
						leftPath = leftPath[1:]
						nodeL1.flag_value.encoded_prefix = compact_encode(leftPath)
						nodeL1.flag_value.value = newValue
						nodeB.branch_value[index] = mpt.InsertRecursive(nodeL1, nil, "")
					} else if len(leftPath) == 0 {
						nodeB.branch_value[16] = newValue
					}
					if len(leftNibble) > 1 {
						nodeL2 := createLeaf()
						index := leftNibble[0]
						leftNibble = leftNibble[1:]
						nodeL2.flag_value.encoded_prefix = compact_encode(leftNibble)
						nodeL2.flag_value.value = currNode.flag_value.value
						nodeB.branch_value[index] = mpt.InsertRecursive(nodeL2, nil, "")
					} else if len(leftNibble) == 0 {
						nodeB.branch_value[16] = currNode.flag_value.value
					}
					nodeE.flag_value.value = mpt.InsertRecursive(nodeB, nil, "")
					hashE := nodeE.hash_node()
					delete(mpt.db, old_hash)
					mpt.db[hashE] = nodeE
					return hashE
				} else { //create branch //create leaves or leaf //put in the branch[16] field //also check if the leaf has just one value
					nodeB := createBranch()

					if len(decodedHexArray) > 0 {
						decodedHexArray = append(decodedHexArray, 16)
						nodeL2 := createLeaf()
						nodeL2.flag_value.encoded_prefix = compact_encode(decodedHexArray[1:])
						nodeL2.flag_value.value = currNode.flag_value.value
						nodeB.branch_value[decodedHexArray[0]] = mpt.InsertRecursive(nodeL2, nil, "")
					} else if len(decodedHexArray) == 0 {
						nodeB.branch_value[16] = currNode.flag_value.value
					}
					if len(leftPath) > 0 {
						nodeL1 := createLeaf()
						leftPath = append(leftPath, 16)
						nodeL1.flag_value.encoded_prefix = compact_encode(leftPath[1:])
						nodeL1.flag_value.value = newValue
						nodeB.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeL1, nil, "")
					} else {
						nodeB.branch_value[16] = newValue
					}
					hashB := nodeB.hash_node()
					mpt.db[hashB] = nodeB
					delete(mpt.db, old_hash)
					return hashB
				}
			} else if nodeType == 0 || nodeType == 1 { // if equal  //if check_if_equal(decodedHexArray, leftPath) {
				if reflect.DeepEqual(decodedHexArray, leftPath) { //insert in branch value place //check if next node is a leaf //if yes convert it into branch //insert this value in branchvalue[16] //check the lenght of the leaf it 1 then create leaf store empty value //lenght of leaf is 0 .... i think it is the same //lenght of leaf is > 1 ....store in leaf
					currNode.flag_value.value = mpt.InsertRecursive(mpt.db[currNode.flag_value.value], nil, newValue)
					hash := currNode.hash_node()
					delete(mpt.db, old_hash)
					mpt.db[hash] = currNode
					return hash
				} else if check_if_common_path_exists(leftPath, decodedHexArray) {
					counter := 0
					for i := 0; i < len(decodedHexArray); i++ {
						if i == len(decodedHexArray) || i == len(leftPath) {
							break
						}
						if decodedHexArray[i] == leftPath[i] {
							counter = counter + 1
						} else {
							break
						}
					}
					commonPath := leftPath[:counter]
					leftPath = leftPath[counter:]
					leftNibble := decodedHexArray[counter:]
					currNode.flag_value.encoded_prefix = compact_encode(commonPath)
					if len(leftNibble) > 0 { //lib
						nodeBranch := createBranch()
						if len(leftPath) > 0 { //create Branch //left path create branch and leaf
							leftPath = append(leftPath, 16)
							nodeLeaf := createLeaf()
							nodeLeaf.flag_value.encoded_prefix = compact_encode(leftPath[1:])
							nodeLeaf.flag_value.value = newValue
							nodeBranch.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeLeaf, nil, "")
						} else {
							nodeBranch.branch_value[16] = newValue
						}
						//check the left nibble size
						if len(leftNibble) == 0 {
							currNode.flag_value.value = mpt.InsertRecursive(mpt.db[currNode.flag_value.value], leftPath, newValue)
							hashExt := currNode.hash_node()
							delete(mpt.db, old_hash)
							mpt.db[hashExt] = currNode
							return hashExt
						} else if len(leftNibble) == 1 {
							nodeBranch.branch_value[leftNibble[0]] = currNode.flag_value.value
							currNode.flag_value.value = mpt.InsertRecursive(nodeBranch, nil, "")
							hashExt := currNode.hash_node()
							delete(mpt.db, old_hash)
							mpt.db[hashExt] = currNode
							return hashExt
						} else if len(leftNibble) > 1 { //create an extension and store the value of that extension in branch => branch hash store in currNode
							nodeExtension := createExtension()
							nodeExtension.flag_value.encoded_prefix = compact_encode(leftNibble[1:])
							nodeExtension.flag_value.value = currNode.flag_value.value
							nodeBranch.branch_value[leftNibble[0]] = mpt.InsertRecursive(nodeExtension, nil, "")
							currNode.flag_value.value = mpt.InsertRecursive(nodeBranch, nil, "")
							hashExt := currNode.hash_node()
							delete(mpt.db, old_hash)
							mpt.db[hashExt] = currNode
							return hashExt
						}
					} else {
						currNode.flag_value.value = mpt.InsertRecursive(mpt.db[currNode.flag_value.value], leftPath, newValue)
						hashExt := currNode.hash_node()
						delete(mpt.db, old_hash)
						mpt.db[hashExt] = currNode
						return hashExt
					}
				} else { //make extension node a branch and one leaf
					nodeBranch := createBranch()
					leftPath = append(leftPath, 16)
					//	fmt.Println("Left path:", leftPath)
					nodeLeaf := createLeaf() //creating a leaf
					nodeLeaf.flag_value.encoded_prefix = compact_encode(leftPath[1:])
					nodeLeaf.flag_value.value = newValue
					nodeBranch.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeLeaf, nil, "")

					if len(decodedHexArray) == 1 {
						nodeBranch.branch_value[decodedHexArray[0]] = currNode.flag_value.value
					} else if len(decodedHexArray) > 1 {
						nodeExt := createExtension()
						nodeExt.flag_value.encoded_prefix = compact_encode(decodedHexArray[1:])
						nodeExt.flag_value.value = currNode.flag_value.value
						nodeBranch.branch_value[decodedHexArray[0]] = mpt.InsertRecursive(nodeExt, nil, "")
					}
					hashBranch := nodeBranch.hash_node()
					mpt.db[hashBranch] = nodeBranch
					delete(mpt.db, old_hash)
					return hashBranch
				}
			}
		} else if check_node_type(currNode) == 1 {
			if currNode.branch_value[leftPath[0]] == "" {
				//store leftPath[0] create a leaf to store the rest
				leftPath = append(leftPath, 16)
				nodeL3 := createLeaf()
				nodeL3.flag_value.encoded_prefix = compact_encode(leftPath[1:])
				nodeL3.flag_value.value = newValue
				currNode.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeL3, nil, "")
			} else if currNode.branch_value[leftPath[0]] != "" {
				index := leftPath[0]
				nextNode := mpt.db[currNode.branch_value[leftPath[0]]]
				leftPath = leftPath[1:]
				currNode.branch_value[index] = mpt.InsertRecursive(nextNode, leftPath, newValue)
			}
			hashBr := currNode.hash_node()
			delete(mpt.db, old_hash)
			mpt.db[hashBr] = currNode
			return hashBr
		}
	}
	return ""
}

//checks if there is common between two slice
func check_if_common_path_exists(leftPath []uint8, decodedHexArray []uint8) bool {
	if len(leftPath) > 0 && len(decodedHexArray) > 0 {
		if leftPath[0] == decodedHexArray[0] {
			return true
		}
	}
	return false
}

//creates branch
func createBranch() Node {
	node := Node{}
	node.node_type = 1
	return node
}

//creates exension
func createExtension() Node {
	node := Node{}
	node.node_type = 2
	return node
}

//creates leaf
func createLeaf() Node {
	node := Node{}
	node.node_type = 2
	return node
}

//Delete function
func (mpt *MerklePatriciaTrie) Delete(key string) string {
	//	fmt.Println("VALUE to be Deleted :", key)
	pathLeft := string_to_hex_array(key)
	currentNode := mpt.db[mpt.Root]
	value, err := mpt.deleteRecursive(mpt.Root, currentNode, pathLeft, []string{})
	if err != nil {
		return value
	}
	return ""
}

func (mpt *MerklePatriciaTrie) deleteRecursive(nodeKey string, currentNode Node, pathLeft []uint8, hashStack []string) (string, error) {
	if len(pathLeft) > 0 && check_node_type(currentNode) != 0 { //path length >0
		if check_node_type(currentNode) == 1 { // if branch check if the next node is a leaf to disconnect the leaf
			if currentNode.branch_value[pathLeft[0]] != "" {
				hash := currentNode.branch_value[pathLeft[0]]
				oldhash := currentNode.hash_node()
				index := pathLeft[0]
				pathLeft = pathLeft[1:]
				next_node := mpt.db[hash]
				if check_node_type(next_node) == 2 {
					hex_prefix_array := ascii_to_hex(next_node.flag_value.encoded_prefix)
					if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) {
						if len(pathLeft) == 0 && len(compact_decode(next_node.flag_value.encoded_prefix)) == 0 {
							currentNode.branch_value[index] = ""
							mpt.db[oldhash] = currentNode
						} else if len(pathLeft) > 0 && len(next_node.flag_value.encoded_prefix) > 0 {
							if reflect.DeepEqual(pathLeft, compact_decode(next_node.flag_value.encoded_prefix)) {
								currentNode.branch_value[index] = ""
								mpt.db[oldhash] = currentNode
							}
						}
					}
				}
				hashStack = append(hashStack, oldhash)
				return mpt.deleteRecursive(hash, mpt.db[hash], pathLeft, hashStack)
			}
		} else if currentNode.node_type == 2 { //ext or leaf and pathleft >0
			hex_prefix_array := ascii_to_hex(currentNode.flag_value.encoded_prefix)
			oldhash := currentNode.hash_node()
			if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension
				nodePath := compact_decode(currentNode.flag_value.encoded_prefix)
				if reflect.DeepEqual(nodePath, pathLeft[:len(nodePath)]) {
					pathLeft = pathLeft[len(nodePath):]
					if len(pathLeft) == 0 { // pathleft is zero now
						hash := currentNode.flag_value.value
						hashStack = append(hashStack, oldhash)
						//pathleft=0, use value from ext and check in branchvalu[16]
						if mpt.db[hash].branch_value[16] != "" { //value found
							nextnode := mpt.db[hash]
							valuetoreturn := nextnode.branch_value[16]
							nextnode.branch_value[16] = ""
							mpt.db[hash] = nextnode //after cheching how many values in branch_value // delete(mpt.db, nodeKey)
							hashStack = append(hashStack, hash)
							mpt.rearrange(hashStack)
							return valuetoreturn, nil
						}
					} else if len(pathLeft) > 0 { //add to hashstack call on next node
						hash := currentNode.flag_value.value
						hashStack = append(hashStack, oldhash) //adding current ext
						return mpt.deleteRecursive(hash, mpt.db[hash], pathLeft, hashStack)
					}
				}
			} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) { //leaf //pathleft >0
				nodePath := compact_decode(currentNode.flag_value.encoded_prefix)
				if reflect.DeepEqual(nodePath, pathLeft) {
					delete(mpt.db, nodeKey) //delete node
					mpt.rearrange(hashStack)
					return currentNode.flag_value.value, nil
				}
			}
		}
	} else if len(pathLeft) == 0 && currentNode.node_type != 0 { //pathlength ==0
		if currentNode.node_type == 1 { // branch
			if currentNode.branch_value[16] != "" {
				previous_hash := currentNode.hash_node()
				previous_value := currentNode.branch_value[16]
				currentNode.branch_value[16] = ""
				mpt.db[previous_hash] = currentNode
				hashStack = append(hashStack, previous_hash) //adding current ext
				mpt.rearrange(hashStack)
				return previous_value, nil
			}
		} else if currentNode.node_type == 2 { //ext or leaf
			hex_empty := ascii_to_hex(currentNode.flag_value.encoded_prefix)
			if hex_empty[0] == 2 || hex_empty[0] == 3 { //leaf
				hex_empty := ascii_to_hex(currentNode.flag_value.encoded_prefix)
				if reflect.DeepEqual(hex_empty, []uint8{2, 0}) {
					value := currentNode.flag_value.value
					delete(mpt.db, currentNode.hash_node())
					mpt.rearrange(hashStack)
					return value, nil
				}
			}
		}
	}
	return "", errors.New("path_not_found")
}

//function which gets the new node, calculates hash and inserts in db
func (mpt *MerklePatriciaTrie) insert_new_node(node Node) string {
	hash := node.hash_node()
	mpt.db[hash] = node
	return hash
}

// function that gives the number of values present in branch
func check_branch_for_values(node Node) (int, int) {
	numValues := 0
	n := 0
	for i := 0; i < 17; i++ {
		if node.branch_value[i] != "" {
			numValues++
			if numValues == 1 {
				n = i
			}
		}
	}
	return numValues, n
}

// rearrangeDeletedTrie rearranges the MPT
func (mpt *MerklePatriciaTrie) rearrange(hashStack []string) {
	counter := len(hashStack) - 1
	mpt.rearrangeRecursive(hashStack, counter, "")
}

func (mpt *MerklePatriciaTrie) rearrangeRecursive(hashStack []string, counter int, currenthash string) bool {
	if counter == -1 { // if the lenght of stack is negative we have traversed till the root, store the new hash
		mpt.Root = currenthash
		return true
	}
	if counter == len(hashStack)-1 { // first loop
		if len(hashStack) == 1 { // only one element in stack
			if check_node_type(mpt.db[hashStack[counter]]) == 1 { //just one node in stack and should be a branch
				currNode := mpt.db[hashStack[counter]]
				numValues, n := check_branch_for_values(currNode)
				if numValues == 1 { //number of values in the branch == 1 check next node
					nextnode := mpt.db[currNode.branch_value[n]]
					nodetype := check_node_type(nextnode)
					if nodetype == 1 { //if branch make the current an extension
						nodeE := createExtension()
						nodeE.flag_value.encoded_prefix = compact_encode([]uint8{uint8(n)})
						delete(mpt.db, hashStack[counter])
						hash := mpt.insert_new_node(nodeE)
						return mpt.rearrangeRecursive(hashStack, counter-1, hash)
					} else if nodetype == 2 {
						_, nodeType := check_if_leaf_or_extension(nextnode)
						decodedHexArray := compact_decode(nextnode.flag_value.encoded_prefix)

						if nodeType == 0 || nodeType == 1 { //if extension club to make an extension together
							nodeE := createExtension()
							combinedHexArray := append([]uint8{uint8(n)}, compact_decode(nextnode.flag_value.encoded_prefix)...)
							nodeE.flag_value.encoded_prefix = compact_encode(combinedHexArray)
							nodeE.flag_value.value = nextnode.flag_value.value
							delete(mpt.db, hashStack[counter])
							delete(mpt.db, hashStack[counter-1])
							hash := mpt.insert_new_node(nodeE)
							return mpt.rearrangeRecursive(hashStack, counter-1, hash)
						} else if nodeType == 2 || nodeType == 3 { //if leaf make a leaf
							nodeL := createLeaf()
							decodedHexArray = append([]uint8{uint8(n)}, decodedHexArray...)
							decodedHexArray = append(decodedHexArray, 16)
							nodeL.flag_value.encoded_prefix = compact_encode(decodedHexArray)
							nodeL.flag_value.value = nextnode.flag_value.value
							delete(mpt.db, hashStack[counter])
							hash := mpt.insert_new_node(nodeL)
							return mpt.rearrangeRecursive(hashStack, counter-1, hash)
						}
					}
				} else if numValues > 1 { // number of values > 1
					node := mpt.db[hashStack[counter]]
					delete(mpt.db, hashStack[counter])
					hash := mpt.insert_new_node(node)
					return mpt.rearrangeRecursive(hashStack, counter-1, hash)
				}
			}
		} else if len(hashStack) > 1 {
			if mpt.db[hashStack[counter]].node_type == 1 { //branch
				currNode := mpt.db[hashStack[counter]]
				numValues, n := check_branch_for_values(currNode)
				if numValues == 1 {
					if n == 16 { //convert to leaf and store value with key = "" //put the value in branch[16] of the following branch
						prevnode := mpt.db[hashStack[counter-1]]
						prevnodetype := check_node_type(prevnode)
						if prevnodetype == 1 { // branch create a empty(key) leaf node with value (branch_value[16])
							nodeL := createLeaf()
							nodeL.flag_value.encoded_prefix = compact_encode([]uint8{16})
							nodeL.flag_value.value = currNode.branch_value[16]
							delete(mpt.db, hashStack[counter])
							hash := mpt.insert_new_node(nodeL)
							return mpt.rearrangeRecursive(hashStack, counter-1, hash)
						} else if prevnodetype == 2 {
							_, nodeType := check_if_leaf_or_extension(prevnode)
							if nodeType == 0 || nodeType == 1 { //extension //add the value to prev Ext node - and make a leaf
								hexArray := compact_decode(prevnode.flag_value.encoded_prefix)
								hexArray = append(hexArray, 16)
								nodeL := createLeaf()
								nodeL.flag_value.encoded_prefix = compact_encode(hexArray)
								nodeL.flag_value.value = currNode.branch_value[16]
								delete(mpt.db, hashStack[counter])
								delete(mpt.db, hashStack[counter-1])
								hash := mpt.insert_new_node(nodeL)
								return mpt.rearrangeRecursive(hashStack, counter-2, hash) //?
							}
						}
					} else if n < 16 && n >= 0 {
						nextnode := mpt.db[currNode.branch_value[n]]
						nodetype := check_node_type(nextnode)
						u := uint8(n)
						if nodetype == 1 {
							//Creating Extension node //check previous node if extension //merger the previous with the current and create extension and update previous // flag.value with value stored at index of current branch
							var newArray []uint8
							var flag int
							previousNode := mpt.db[hashStack[counter-1]]
							if check_node_type(previousNode) == 2 {
								_, nodeType := check_if_leaf_or_extension(previousNode)
								if nodeType == 0 || nodeType == 1 {
									flag = 1
									newArray = compact_decode(previousNode.flag_value.encoded_prefix)
								}
							}
							newArray = append(newArray, []uint8{u}...)
							nodeE := createExtension()
							nodeE.flag_value.encoded_prefix = compact_encode(newArray)
							nodeE.flag_value.value = currNode.branch_value[n]
							delete(mpt.db, hashStack[counter])
							hash := mpt.insert_new_node(nodeE)
							if flag == 1 {
								delete(mpt.db, hashStack[counter-1])
								return mpt.rearrangeRecursive(hashStack, counter-2, hash)
							} else {
								return mpt.rearrangeRecursive(hashStack, counter-1, hash)
							}
						} else if nodetype == 2 {
							_, nodeType := check_if_leaf_or_extension(nextnode)
							if nodeType == 0 || nodeType == 1 { //extension
								var newArray []uint8
								previousNode := mpt.db[hashStack[counter-1]]
								var flag int
								if previousNode.node_type == 2 {
									_, nodeType := check_if_leaf_or_extension(previousNode)
									if nodeType == 0 || nodeType == 1 {
										flag = 1
										newArray = compact_decode(previousNode.flag_value.encoded_prefix)
									}
								}
								newArray = append(newArray, []uint8{u}...)
								hexArray := compact_decode(nextnode.flag_value.encoded_prefix)
								newArray = append(newArray, hexArray...)
								nextHex := currNode.branch_value[n]
								nextnode.flag_value.encoded_prefix = compact_encode(newArray)
								delete(mpt.db, hashStack[counter])
								delete(mpt.db, nextHex)
								hash := mpt.insert_new_node(nextnode)
								if flag == 1 {
									delete(mpt.db, hashStack[counter-1])
									return mpt.rearrangeRecursive(hashStack, counter-2, hash)
								} else {
									return mpt.rearrangeRecursive(hashStack, counter-1, hash)
								}
							} else if nodeType == 2 || nodeType == 3 { //leaf
								//check if the previous is extension // if so club it with the leaf
								var flag int
								var hexArray []uint8
								previousnode := mpt.db[hashStack[counter-1]]
								if check_node_type(previousnode) == 2 {
									_, nodeType := check_if_leaf_or_extension(previousnode)
									if nodeType == 0 || nodeType == 1 {
										flag = 1
										hexArray = compact_decode(previousnode.flag_value.encoded_prefix)
									}
								}
								hexArray = append(hexArray, []uint8{u}...)
								toBeAddedArray := compact_decode(nextnode.flag_value.encoded_prefix)
								hexArray = append(hexArray, toBeAddedArray...)
								hexArray = append(hexArray, 16)
								nextnode.flag_value.encoded_prefix = compact_encode(hexArray)
								delete(mpt.db, hashStack[counter])
								hash := mpt.insert_new_node(nextnode)
								if flag == 1 {
									delete(mpt.db, hashStack[counter-1])
									return mpt.rearrangeRecursive(hashStack, counter-2, hash)
								} else {
									return mpt.rearrangeRecursive(hashStack, counter-1, hash)
								}
							}
						}
					}
				} else if numValues > 1 { //if number of values in the branch is > 1 just update the hzsh in the db
					node := mpt.db[hashStack[counter]]
					delete(mpt.db, hashStack[counter])
					counter = counter - 1
					hash := mpt.insert_new_node(node)
					return mpt.rearrangeRecursive(hashStack, counter, hash)
				}
			}
		}
	} else { // for the later nodes just update the hash if they are branch or extension
		node := mpt.db[hashStack[counter]]
		if check_node_type(node) == 1 {
			for i := 0; i < 17; i++ {
				if node.branch_value[i] == hashStack[counter+1] {
					node.branch_value[i] = currenthash
				}
			}
			delete(mpt.db, hashStack[counter])
			hash := mpt.insert_new_node(node)
			return mpt.rearrangeRecursive(hashStack, counter-1, hash)
		} else if check_node_type(node) == 2 {
			node.flag_value.value = currenthash
			delete(mpt.db, hashStack[counter])
			hash := mpt.insert_new_node(node)
			return mpt.rearrangeRecursive(hashStack, counter-1, hash)
		}
	}
	return false
}

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value + string(node.flag_value.encoded_prefix)
	}
	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) Initial() {
	// mpt.db = map[string]Node{}
	mpt.db = make(map[string]Node)
	mpt.Root = ""
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0]/16 < 2
}

func TestCompact() {
	test_compact_encode()
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.Root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

//Order_nodes method orders the nodes
func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart"+cur_hash+"HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart"+cur_hash+"HashEnd", fmt.Sprintf("Hash%v", i), -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}

//GetAll method Project2:Gets the key-value pairs in the mpt in map
func (mpt *MerklePatriciaTrie) GetAll() map[string]string {
	keyValueMap := make(map[string]string)
	root := mpt.Root
	rootNode := mpt.db[root]
	path := []uint8{}
	mpt.GetAllRecursive(rootNode, path, keyValueMap)
	return keyValueMap
}

//GetAllRecursive method Project2:Gets the key-value pairs in the mpt in map
func (mpt *MerklePatriciaTrie) GetAllRecursive(node Node, path []uint8, keyValueMap map[string]string) {
	if node.node_type == 1 { //branch
		for i := 0; i < 17; i++ {
			if node.branch_value[i] != "" {
				if i == 16 {
					keyValueMap[hexToString(path)] = node.branch_value[16]
				} else {
					newPath := append(path, []uint8{uint8(i)}...)
					nextNode := mpt.db[node.branch_value[i]]
					mpt.GetAllRecursive(nextNode, newPath, keyValueMap)
				}
			}
		}
	} else if node.node_type == 2 {
		_, nodeType := check_if_leaf_or_extension(node)
		if nodeType == 0 || nodeType == 1 { //extension
			decodedHexArray := compact_decode(node.flag_value.encoded_prefix)
			path = append(path, decodedHexArray...)
			nextNode := mpt.db[node.flag_value.value]
			mpt.GetAllRecursive(nextNode, path, keyValueMap)
		} else if nodeType == 2 || nodeType == 3 { //leaf
			decodedHexArray := compact_decode(node.flag_value.encoded_prefix)
			if len(decodedHexArray) > 0 {
				path = append(path, decodedHexArray...)
			}
			keyValueMap[hexToString(path)] = node.flag_value.value
		}
	}
}

func main() {
	mpt := &MerklePatriciaTrie{make(map[string]Node), ""}
	fmt.Println("inserting do")
	mpt.Insert("do", "verb")
	fmt.Println("inserting dog")
	mpt.Insert("dog", "puppy")
	fmt.Println("inserting doge")
	mpt.Insert("doge", "coin")
	fmt.Println("inserting horse")
	mpt.Insert("horse", "stallion")
	mpt.GetAll()
}
