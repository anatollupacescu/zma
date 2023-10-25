package bmt

import (
	"testing"
)

func NewTestBmt() *Bmt[string] {
	return New[string](func(l, r string) string {
		return l + r
	})
}

func TestProof(t *testing.T) {
	t.Run("two files, proof len is 1", func(t *testing.T) {
		bmt := NewTestBmt()
		_ = bmt.Add("1test")
		_ = bmt.Add("2test")

		proofs := bmt.Proof(0)
		if len(proofs) != 1 {
			t.Fatal("expected one proof")
		}

		if bmt.Size() != 3 {
			t.Fatal("want size 3")
		}
	})

	t.Run("three files, proofs contain zero hash", func(t *testing.T) {
		bmt := NewTestBmt()
		_ = bmt.Add("a")
		_ = bmt.Add("b")
		_ = bmt.Add("c")

		proofs := bmt.Proof(2)
		if len(proofs) != 2 {
			t.Fatal("expected one proof")
		}

		if proofs[0] != "" {
			t.Fatal("expected zero hash")
		}
	})

	t.Run("four files, proof len is 2", func(t *testing.T) {
		bmt := NewTestBmt()
		_ = bmt.Add("a")
		_ = bmt.Add("b")
		_ = bmt.Add("c")
		_ = bmt.Add("d")

		proofs := bmt.Proof(0)
		if len(proofs) != 2 {
			t.Fatal("expected one proof")
		}

		if proofs[0] != "b" {
			t.Fatal("expected b")
		}

		if proofs[1] != "cd" {
			t.Fatal("expected cd")
		}
	})
}

func TestAdd(t *testing.T) {
	t.Run("singe file hash is root hash", func(t *testing.T) {
		bmt := NewTestBmt()
		got := bmt.Add("test test")
		want := bmt.RootSum()

		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("two file hashes geneate a new root", func(t *testing.T) {
		bmt := NewTestBmt()
		original := bmt.Add("test test")
		_ = bmt.Add("test test")

		sum := bmt.RootSum()

		if original == sum {
			t.Fatal("hash sums should not be equal")
		}

		if len(bmt.sums) != 2 {
			t.Fatal("want two base hashes")
		}

		if len(bmt.next.sums) != 1 {
			t.Fatal("want one root hash")
		}
	})

	t.Run("four files combine hash correctly", func(t *testing.T) {
		bmt := NewTestBmt()
		_ = bmt.Add("a")
		_ = bmt.Add("b")
		_ = bmt.Add("c")
		_ = bmt.Add("d")

		p0, p1 := bmt.next.sums[0], bmt.next.sums[1]
		if p0 != "ab" || p1 != "cd" {
			t.Fatalf("wrong combinations: %s", bmt.next.sums)
		}

		if bmt.next.next.sums[0] != "abcd" {
			t.Fatalf("wrong root hash: %s", bmt.next.next.sums[0])
		}

		if bmt.Size() != 7 {
			t.Fatalf("want size 7, got %d", bmt.Size())
		}
	})
}
