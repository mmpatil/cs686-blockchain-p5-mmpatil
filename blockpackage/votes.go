package blockpackage

import (
	p1 "../p1"
	p5 "../p5/data"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type VotesNotFinalized struct {
	Votes map[string]p5.RequestResponse
	mux   sync.Mutex
}

type FinalizedVotes struct {
	FinalizedVotes   map[string]string `json:"FinalizedVotes"`
	TotalVotes       int               `json:"TotalVotes"`
	CandidateVoteMap map[int]int       `json:"CandidateVoteMap"`
	mux              sync.Mutex
}

func InitializieVoteMaps() (VotesNotFinalized, FinalizedVotes) {

	notfinalizedVotes := VotesNotFinalized{
		Votes: make(map[string]p5.RequestResponse),
	}

	finalizedVotes := FinalizedVotes{
		FinalizedVotes:   make(map[string]string),
		TotalVotes:       0,
		CandidateVoteMap: make(map[int]int),
	}

	finalizedVotes.CandidateVoteMap[1] = 0
	finalizedVotes.CandidateVoteMap[2] = 0

	return notfinalizedVotes, finalizedVotes
}

func (finalizedVotes *FinalizedVotes) IfValidBlock(newBlock Block) bool {
	fmt.Println("In votes.go IfValidBlock")
	finalizedVotes.mux.Lock()
	defer finalizedVotes.mux.Unlock()
	mpt := p1.MerklePatriciaTrie{}
	mpt = newBlock.Value
	mptMap := mpt.GetAll()

	//var finalizedVotesNew FinalizedVotes
	//finalizedVotesNew  = newBlock.Header.FinalizedVotesStruct
	//
	//if finalizedVotes.TotalVotes >= finalizedVotesNew.TotalVotes{
	//	for k,_ := range finalizedVotesNew.CandidateVoteMap{
	//		if finalizedVotesNew.CandidateVoteMap[k] == finalizedVotes.CandidateVoteMap[k]{
	//			return false
	//		}
	//	}
	//	return false
	//}

	for k, _ := range mptMap {
		_, exists := finalizedVotes.FinalizedVotes[k]
		if exists == true {
			fmt.Println("Result:Wrong BlockBeat")
			return false
		}
		//else{
		//	finalizedVotes.FinalizedVotes[k] = v
		//	finalizedVotes.TotalVotes += 1
		//}
	}
	//for key, value := range mptMap {
	//	finalizedVotes.FinalizedVotes[key] = value
	//	finalizedVotes.TotalVotes += 1
	//	fmt.Println("finalizedVotes",finalizedVotes.FinalizedVotes)
	//}
	fmt.Println("Result:Valid BlockBeat")
	return true
}

func (finalizedVotes *FinalizedVotes) InsertInToFinalizedVotes(newBlock Block) {
	fmt.Println("----------------------------------InsertInToFinalizedVotes-----------------")
	finalizedVotes.mux.Lock()
	defer finalizedVotes.mux.Unlock()

	mpt := p1.MerklePatriciaTrie{}
	mpt = newBlock.Value
	mptMap := mpt.GetAll()
	for key, value := range mptMap {
		reqresp := p5.RequestResponse{}
		finalizedVotes.FinalizedVotes[key] = value
		err := json.Unmarshal([]byte(value), &reqresp)
		if err != nil {
			log.Fatal("Error while counting candidate so unmarshal")
		}
		id := reqresp.User.CandidateId
		value, exists := finalizedVotes.CandidateVoteMap[id]
		if exists {
			finalizedVotes.TotalVotes += 1
			value++
			finalizedVotes.CandidateVoteMap[id] = value
			finalizedVotes = &newBlock.Header.FinalizedVotesStruct //justnow p5
		}
		fmt.Println("*******************************************************!!!!!!!!!!!!!!!!!!!!!!finalizedVotes", finalizedVotes.TotalVotes)
	}
}

func (finalizedVotes *FinalizedVotes) checkIfVoteInBlockchain(tempPublickey string) bool {
	finalizedVotes.mux.Lock()
	defer finalizedVotes.mux.Unlock()

	_, exists := finalizedVotes.FinalizedVotes[tempPublickey]
	if exists == true {
		return true
	}
	return false
}

func PrepareMPT(finalizedVotes FinalizedVotes, votesNotFinalized VotesNotFinalized, initial bool) (p1.MerklePatriciaTrie, bool) {
	if initial {
		mpt := p1.MerklePatriciaTrie{}
		mpt.Initial()
		return mpt, true
	}
	//finalizedVotes.mux.Lock()
	//defer finalizedVotes.mux.Unlock()
	//votesNotFinalized.mux.Lock()
	//defer votesNotFinalized.mux.Unlock()
	var valid bool
	valid = false
	mpt := p1.MerklePatriciaTrie{}
	mpt.Initial()
	var arrayPK []string
	notfinalizedvotes := votesNotFinalized.Votes
	//fmt.Println("Lenght of votesNotFinalized.Votes:",len(notfinalizedvotes))
	if len(notfinalizedvotes) > 0 {
		fmt.Println("len(notfinalizedvotes) is:", len(notfinalizedvotes))
		for k, v := range notfinalizedvotes {
			_, exists := finalizedVotes.FinalizedVotes[k]
			if exists == false {
				//if finalizedVotes.checkIfVoteInBlockchain(k) == false {
				reqResString := v.EncodeRequestRespToJson()
				mpt.Insert(k, reqResString)
				arrayPK = append(arrayPK, k)
				valid = true
				fmt.Println("MPT is being constructed")
			}
			//	delete(votedNotFinalized, k)
		}
		for _, value := range arrayPK {
			delete(votesNotFinalized.Votes, value)
		}
	}
	return mpt, valid
}

func (finalizedVotes *FinalizedVotes) ExistsInFinalizedVote(tempPublicKey string) bool {
	finalizedVotes.mux.Lock()
	defer finalizedVotes.mux.Unlock()

	_, exists := finalizedVotes.FinalizedVotes[tempPublicKey]
	return exists
}

func (votesNotFinalized *VotesNotFinalized) ExistsInVotesNotFinalized(tempPublicKey string) bool {
	votesNotFinalized.mux.Lock()
	defer votesNotFinalized.mux.Unlock()

	_, exists := votesNotFinalized.Votes[tempPublicKey]
	return exists
}

func (votesNotFinalized *VotesNotFinalized) InsertVotedNotFinalized(publicKey string, reqresp p5.RequestResponse) bool {
	votesNotFinalized.mux.Lock()
	defer votesNotFinalized.mux.Unlock()
	//if votesNotFinalized.ExistsInVotesNotFinalized(publicKey) == false {
	_, exists := votesNotFinalized.Votes[publicKey]
	if exists == false {
		votesNotFinalized.Votes[publicKey] = reqresp
		return true
	}
	return false
}

func (finalizedVotes *FinalizedVotes) InsertFinalizedVotes(publicKey string, jsonString string) bool {
	finalizedVotes.mux.Lock()
	defer finalizedVotes.mux.Unlock()

	if finalizedVotes.ExistsInFinalizedVote(publicKey) == false {
		finalizedVotes.FinalizedVotes[publicKey] = jsonString
		return true
	}
	return false
}
