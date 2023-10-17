package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/anatollupacescu/zma/bmt"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("usage: hasher path/to/dir")
	}

	dir, err := filepath.Abs(filepath.Clean(args[1]))
	if err != nil {
		log.Fatalf("absolute path %s: %v", dir, err)
		return
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("read dir %s: %v", dir, err)
	}

	cbmt := bmt.New()
	var rootSum [32]byte

	for _, file := range files {
		fileLoc := filepath.Join(dir, file.Name())
		data, err := os.Open(fileLoc)
		if err != nil {
			log.Fatalf("open %s: %v", fileLoc, err)
			return
		}

		buf, err := io.ReadAll(data)
		if err != nil {
			log.Fatalf("read contents: %v", err)
			return
		}

		rootSum = cbmt.Add(buf)
	}

	fmt.Println(hex.EncodeToString(rootSum[:]))
}
