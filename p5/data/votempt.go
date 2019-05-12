package data

import (
	"../../p1"
	"crypto/rsa"
)

type voteMPT struct {
	mpt       p1.MerklePatriciaTrie `json:"mpt"`
	signature []byte                `json:"signature"`
	publicKey *rsa.PublicKey        `json:"publicKey"`
}
