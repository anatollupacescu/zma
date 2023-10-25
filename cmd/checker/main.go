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
	Sum string `json:"sum"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: checker -root $YOUR_ROOT_HASH -proofs $PROOFS_ARRAY -index $INDEX filename.png")
	}

	filename := os.Args[7]

	rootSum := flag.String("root", "", "root hash sum")
	proofs := flag.String("proofs", "[]", "the set of proof sums")
	index := flag.Int("index", 0, "index of the file")
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

		var left bool
		switch *index {
		case 0:
			left = true
		case 1:
		default:
			if *index%2 == 0 {
				left = true
			}
		}

		if left {
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
