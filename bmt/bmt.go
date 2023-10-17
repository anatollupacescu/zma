package bmt

import (
	"golang.org/x/crypto/sha3"
)

type Bmt struct {
	root [32]byte
	node *node

	addF func([32]byte) [32]byte
}

type node struct {
	sums [][32]byte
	next *node
}

func New() *Bmt {
	b := &Bmt{
		node: new(node),
	}
	b.addF = b.initAdd
	return b
}

func (b *Bmt) RootSum() [32]byte {
	return b.root
}

func (b *Bmt) initAdd(in [32]byte) [32]byte {
	b.node.sums = append(b.node.sums, in)
	b.node.next = new(node)
	b.node.next.sums = append(b.node.next.sums, in)
	b.addF = func(sum [32]byte) [32]byte {
		return b.node.add(sum, false)
	}
	return in
}

func (b *Bmt) Len() int {
	if b.node.next == nil {
		return 0
	}

	return len(b.node.sums)
}

func (b *Bmt) Add(file []byte) [32]byte {
	sum := sha3.Sum256(file)
	b.root = b.addF(sum)
	return b.root
}

var zero [32]byte

func (n *node) add(sum [32]byte, replace bool) (digest [32]byte) {
	if replace {
		n.sums = n.sums[:len(n.sums)-1]
		replace = false
	}

	n.sums = append(n.sums, sum)

	if len(n.sums) == 1 {
		return n.sums[0]
	}

	var left, right [32]byte

	length := len(n.sums)

	if length%2 == 0 {
		left = n.sums[length-2]
		right = n.sums[length-1]
		replace = true
	} else {
		left = n.sums[length-1]
		right = zero
	}

	combined := sha3.Sum256(append(left[:], right[:]...))

	if n.next == nil {
		n.next = new(node)
		replace = false
	}

	return n.next.add(combined, replace)
}

type Proof struct {
	Left bool
	Sum  [32]byte
}

func (b *Bmt) Proof(i int) []Proof {
	return b.node.proof(i, nil)
}

func (n *node) proof(index int, acc []Proof) []Proof {
	var (
		proofIndex int
		left       bool
	)

	if index == 0 {
		proofIndex = 1
		left = true
	} else if index == 1 {
		proofIndex = 0
	} else if index%2 == 0 {
		proofIndex = index + 1
		left = true
	} else {
		proofIndex = index - 1
	}

	if proofIndex < len(n.sums) {
		acc = append(acc, Proof{
			Sum:  n.sums[proofIndex],
			Left: left,
		})
	} else {
		acc = append(acc, Proof{
			Sum:  zero,
			Left: left,
		})
	}

	if n.isLast() {
		return acc
	}

	if left {
		index++
	}

	index /= 2

	return n.next.proof(index, acc)
}

func (n *node) isLast() bool {
	return n.next != nil && n.next.next == nil
}
