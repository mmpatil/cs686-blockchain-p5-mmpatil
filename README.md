# cs686_BlockChain_P5


# Transparent Voting System 
This is a three part application having a user interface for voters to interact with the voting system(For registration /voting and checking the blockchain for their vote). An Registration server which verifies the identity and And only the the Voting page is displayet to the User. And a distributed system to count the votes and keep the process of voting completely visible. (Blockchain is used as a blackbox)

Step 1: Register with NationaID

Step 2: Sigin with valid Public-Private Key Pair to get the Voting Page
	
Step 3: Voting Page displayed for valid user, get the peerlist from 6686 to send a vote to a live listener

Step 4: Blockchain takes care of sending it to other peers update the count for the candidates

Step 5: Can check the vote count of candidate

Not Completed: Allowing the user to fetch its vote using public key from a live listener(peerList).

Ensures that the your vote is counted, i.e. it is a part of the CannonicalBlockchain.

Provides Transparency and Anonymity.

## Votes.go

Has two structures :

-VotesNotFinalized
	map Votes (key:pubicKey value:RequestResponse)
	mux sync.Mutex lock

-FinalizedVotes
	map FinalizedVotes (key: PublicKey value:RequestResponse in json string)
	int TotalVotes
	map CandidateVoteMap (key : CandidateId value:number of votes)
	mux syn.Mutex lock

Functions:
-InitializieVoteMaps() :Initializes the FinalizedVotes and VotesNotFinalized struct

-IfValidBlock() :IfValidBlock method checks Block is valid if the MPT does not contain any public-key that exists in 			FinalizedVotes/ blockchain

-InsertInToFinalizedVotes() :This inserts the mpt of a block into the finalizedVotes struct of a peer and updates totalvotes and increases the vote for the corresponding Candidate.

-checkIfVoteInBlockchain() : Checks if the publicKey of a voter already exists in the Blockchain (in other words if a votes has already voted)

-PrepareMPT() : Prepares an mpt by taking votes from the votesNotFinalized pool of votes that are not still counted and valid.

-ExistsInFinalizedVote() : Checks if a publicKey exists in the FinalizedVote struct of a peer.

-ExistsInVotesNotFinalized : Checks if a publicKey exists in the VotesNotFinalized struct of a peer.

-InsertVotedNotFinalized() : Insert a vote into the vote pool or VotedNotFinalized struct.

-InsertFinalizedVotes() : Insert a vote into the FinalizedVotes struct.

## user.go 

Has two structures:
- User
	NationalId
	PublicKey
	PrivateKey
	CandidateId
	
-RequestResponse
	User
	Signature
	PublicKey

-EncodeToJson():Marshall User object

-EncodeRequestRespToJson(): Marshall RequestResponse object

-Vote(): Allows user to vote.

## UserList.go

Has one structure:
UsersList struct:
	UserMap map  (key : publicKey value:User)
	PKMap   map   (key : publicKey value: nationalId)
	mux     sync.Mutex 

-NewUserList(): Initializes newUserList
-Verify(): Verify if the user already exists
-CopyUsersMap() :Copies the usersMap
-CopyPKMap() :Copies the PublicKey map

## routes.go

/startClient - StartClient() : Starts a client/user/voter with initial configuration

/startReg - StartRegistrationServer() :Starts a registration server with initial configuration

/signup - SignUp(): Display SignUp page

/signin - SignIn(): Display SignIn page

/registerClient - RegisterClient(): Registers a Client with the registration Server/ Authentication Server.

/checkUser - CheckUser(): When a voter/user signs in its Public-Key and Private-Key pair is verified with the registration 		server.

/check - Check() : Checks if the new user/voter trying to register already exists

/voteDetails - VoteDetails(): Gives the user/voter the details of its vote

/vote - Vote(): Vote sent to the miner

/getPeerList - GetPeerList():User/Voter fetches a peerlist from the peer 6686

/showMPT - ShowMPT(): Display MPT

/showBlockAtHeight/{height} - ShowBlockAtHeight(): Shows block at a particular height at a miner

/showVoteUser - ShowVoteUser() Shows the vote details of a user to that user.

/clientVote - ClientVote(): Allows a valid user/voter to vote

## block.go
	The header of the Block and the BlockJson will have an additional field FinalizedVotesStruct which each miner will send each other after updating.

## PEERS (BLOCKCHAIN)

Additional map
VotesInBlockchain := map(key = public_key, value = VoteMPT)
CandidateAndVoteCountMap := map(key = candidate_id , value = count)

### 1)Receive(w,r):
Get the VoteMPT from the body.
Verify the signature in the VoteMPT
If the signature is valid 
	mpt = voteMPT.mpt 
	Verify(mpt) //verify for valid signatures for each key-value in mpt
	Check if any public_key(key) in mpt exists in the VotesInBlockchain.
	If mpt is valid(i.e. the public_keys in mpt do not exist in the VotesInBlockchain) 
	Insert in VotesInBlockchain each key value pair in mpt
	In CandidateAndVoteCountMap increment the count for each vote for respective 			candidate
	PrepareHeartBeat(mpt, CandidateAndVoteCountMap);

### 2)RecieveHeartBeat() of Each Peer:
Get the VoteMPT from the body.
Verify the signature in the VoteMPT
If the signature is valid 
	mpt = voteMPT.mpt 
	Verify(mpt) //verify for valid signatures for each key-value in mpt
	Check if any public_key(key) in mpt exists in the VotesInBlockchain.
	If mpt is valid(i.e. the public_keys in mpt do not exist in the VotesInBlockchain) 
	Insert in VotesInBlockchain each key value pair in mpt
	In CandidateAndVoteCountMap increment the count for each vote for respective 			candidate
	ForwardHeartBeat()


------------------------------------------------------------------------------------------------------------
# Crypto

## functions ->
1) create pub - priv key pair
    :using rsa.GenerateKey() method from the cryto package.
    
2) create Signature : creating signature for the message with the private key with SHA256 hash function.

3) verify signature
4) encrypt message
5) decrypt message
