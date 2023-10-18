package kbmt

import (
	"github.com/anatollupacescu/zma/bmt"
	"github.com/ethereum/go-ethereum/crypto"
)

// NewKeccak256Bmt instantiate a new bmt using the go-ethereum's Keccak256
func NewKeccak256Bmt() *bmt.Bmt[[]byte] {
	return bmt.New[[]byte](Comb)
}

func Comb(l, r []byte) []byte {
	return []byte(crypto.Keccak256(append(l, r...)))
}

func Sum(file []byte) []byte {
	return []byte(crypto.Keccak256(file))
}
