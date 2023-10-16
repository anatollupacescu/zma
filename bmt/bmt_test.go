package bmt

import (
	"bytes"
	"testing"
)

func TestProof(t *testing.T) {
	t.Run("two files, proof len is 1", func(t *testing.T) {
		bmt := New()
		_ = bmt.Add([]byte("1test"))
		_ = bmt.Add([]byte("2test"))

		proofs := bmt.Proof(0)
		if len(proofs) != 1 {
			t.Fatal("expected one proof")
		}
	})

	t.Run("three files, proof len is 2", func(t *testing.T) {
		bmt := New()
		_ = bmt.Add([]byte("1test"))
		_ = bmt.Add([]byte("2test"))
		_ = bmt.Add([]byte("3test"))

		proofs := bmt.Proof(1)
		if len(proofs) != 2 {
			t.Fatal("expected one proof")
		}
	})

	t.Run("four files, proof len is 2", func(t *testing.T) {
		bmt := New()
		_ = bmt.Add([]byte("1test"))
		_ = bmt.Add([]byte("2test"))
		_ = bmt.Add([]byte("3test"))
		_ = bmt.Add([]byte("4test"))

		proofs := bmt.Proof(1)
		if len(proofs) != 2 {
			t.Fatal("expected one proof")
		}
	})
}

func TestAdd(t *testing.T) {
	t.Run("singe file hash is root hash", func(t *testing.T) {
		bmt := New()
		buf := []byte("test test")
		got := bmt.Add(buf)

		want := bmt.RootSum()

		if !bytes.Equal(got[:], want[:]) {
			t.Fatal("hash sums not equal")
		}
	})

	t.Run("two file hashes geneate a new root", func(t *testing.T) {
		bmt := New()
		buf := []byte("test test")

		original := bmt.Add(buf)
		_ = bmt.Add(buf)

		want := bmt.RootSum()

		if bytes.Equal(original[:], want[:]) {
			t.Fatal("hash sums should not be equal")
		}

		if len(bmt.node.sums) != 2 {
			t.Fatal("want two base hashes")
		}

		if len(bmt.node.next.sums) != 1 {
			t.Fatal("want one root hash")
		}
	})
}
