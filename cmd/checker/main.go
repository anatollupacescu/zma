package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"os"

	kbmt "github.com/anatollupacescu/zma/keccak256bmt"
)

type proof struct {
	Left bool   `json:"left"`
	Sum  string `json:"sum"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: checker -root $YOUR_ROOT_HASH -proofs $PROOFS_ARRAY filename.png")
	}

	filename := os.Args[5]

	rootSum := flag.String("root", "", "root hash sum")
	proofs := flag.String("proofs", "[]", "the set of proof sums")
	flag.Parse()

	if *rootSum == "" {
		log.Fatal("root sum not provided")
	}

	if *proofs == "[]" {
		log.Fatal("set of proofs not provided")
	}

	var proofArr []proof
	if err := json.Unmarshal([]byte(*proofs), &proofArr); err != nil {
		log.Fatal("invalid json array")
	}

	root, err := hex.DecodeString(*rootSum)
	if err != nil {
		log.Fatal("invalid json array")
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("read file", err)
	}

	contentHashSum := kbmt.Sum(content)

	var sumToMatch = contentHashSum

	for _, proof := range proofArr {
		proofSum, err := hex.DecodeString(proof.Sum)
		if err != nil {
			log.Fatal("decode sum", err)
		}

		if proof.Left {
			sumToMatch = kbmt.Comb(sumToMatch, proofSum)
			continue
		}
		sumToMatch = kbmt.Comb(proofSum, sumToMatch)
	}

	if !bytes.Equal(sumToMatch, root) {
		log.Fatal("hash sums do NOT match")
	}

	log.Println("hash sums match")
}
