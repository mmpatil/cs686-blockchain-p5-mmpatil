# cs686_BlockChain_P5


					###Voting System

This is a three part application having a user interface for voters to interact with the voting system(For registration /voting and checking the blockchain for their vote). An authentication server which verifies the identity and allows a valid voter to vote. And a distributed system to count the votes and keep the process of voting completely visible. 


Datastructures:

##User:
Each user will have a national_id, public_key ,private_key and a candidate_id (to represent whom did the user vote to).

##2. Response
Each Response will have user object , signature and public key.

Required function for User:

1)Register(authentication_server_address)
Get request to the server address with the national id of the user in the body of the  request.
The response will contain the public and the private key pair. (User object with public and private key values set)
Set the public_key and private_key values of this user.

2) Vote(candidate_id, authentication_server_address)
Set the candidate_id for this user indicating the user is voting this candidate for the election (user.candidate_id = candidate_id). 
Create a User(userDummy) object with public_key , candidate_id set.
Create a Response object 
response.User = userDummy
response.signature = Create signature (userDummy,private_key)
Response. public_key = public_key

Post request to the authentication_server_address with the user details in the body. (response object)

3) GetVoteDetails(blockchain_server_address) 
Create a userDummy with public_key value set.
Create a Response object (response)
response.User = userDummy
response.signature = Create signature (userDummy,private_key)
Get  request to the blockchain_server_address (one of the live-peers) with the response object in the body.
Display the result.

3. VoteMPT
	Each VoteMPT will have 
Mpt := MerkelPatriciaTrie (key = public_key, value = jsonstring)
Signature := signature of the authentication server
Public_key : = public_key of the authentication server

3.  Authentication_Server
Each authentication server will have 
public_key, 
private_key, 
PeerList [] (list of peers online), 
UserDetailsMap  := map(key = public_key, value = User object),
VoteCount := integer , number of people voted.
VotedMap := map(key= int, value = VoteResponse)

Required Functions:

1)Register(w,r)
Verify the national_id in the body of the request.
If valid national_id create public-private key-pair
Create a User object assign the national_id, public_key and private_key.
Insert(public_key, User) in UserDetailsMap 
Create a Response object. (registerResponse)
registerResponse.user = User
registerResponse.signature = Create a signature(User,private_key)
resgisterResponse.public_key = public_key
Send a response containing registerResponse object.

2)Vote(w,r)
Verify the signature in the post request.
If the signature is valid 
	Update the user object for the public_key in the UserDetailsMap (candidate_id 			field).
	Increment the VoteCount field (Authentication_server. VoteCount ++)
	Insert(voted, responseObject) in VotedMap
	Send Response(200OK)
Else
	Send Error


3)CountVotes()
Initialise a map VotesCountedMap (key = int, value = public_key)
Count := 1
while(len(VotesCountedMap) < len(VotedMap)){
	mpt , VotesCountedMap, count := PrepareMPT(VotesCountedMap, VotedMap)
	Create VoteMPT object (voteMpt)
	voteMpt.mpt = mpt
	voteMpt.signature = PrepareSignature(mpt,private_key)
	voteMpt.public_key = public_key
	
	send a post request to one of the live peers with voteMpt in the body.
	Sleep(20 seconds)
	Refresh the peerList
}

Print “Counting of votes is done”!!
}	

4)PrepareMPT(VotesCountedMap, VotedMap,int count)
Create a MerklePatriciaTrie object (mpt)
For i= count ; I < count + 10; count++{
	VoteResponse := VotedMap.get(count)
	count++;
	Mpt.Insert(VoteResponse.public_key, json.Marshall(VoteResponse))
	insert in VotesCountedMap(count,VoteResponse.public_key)
}
Return mpt,VotesCountedMap, count



PEERS (BLOCKCHAIN)

Additional map
VotesInBlockchain := map(key = public_key, value = VoteMPT)
CandidateAndVoteCountMap := map(key = candidate_id , value = count)

1)Receive(w,r)
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

2)RecieveHeartBeat() of Each Peer
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



# Crypto

## functions ->
1) create pub - priv key pair
    :using rsa.GenerateKey() method from the cryto package.
    
    
    
2) create Signature : creating signature for the message with the private key with SHA256 hash function.

3) verify signature
4) encrypt message
5) decrypt message
