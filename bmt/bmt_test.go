package bmt_test

import (
	"bytes"
	"testing"

	"github.com/anatollupacescu/zma/bmt"
)

func TestBmt(t *testing.T) {
	t.Run("singe file hash is root hash", func(t *testing.T) {
		bmt := bmt.New()
		buf := []byte("test test")
		got := bmt.Add(buf)

		want := bmt.RootSum()

		if !bytes.Equal(got[:], want[:]) {
			t.Fatal("hash sums not equal")
		}
	})
}
