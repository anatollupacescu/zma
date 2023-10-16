package bmt

import "golang.org/x/crypto/sha3"

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

func (b *Bmt) Proof(i int) [][32]byte {
	return b.node.proof(i, nil)
}

func (n *node) proof(i int, acc [][32]byte) [][32]byte {
	var pfi int

	if i == 0 {
		pfi = 1
	} else if i == 1 {
		pfi = 0
	} else if i%2 == 0 {
		pfi = i - 1
	} else {
		pfi = i + 1
	}

	acc = append(acc, n.sums[pfi])

	if n.isLast() {
		return acc
	}

	isLeft := pfi > i
	if isLeft {
		i++
	}

	i /= 2

	return n.next.proof(i, acc)
}

func (n *node) isLast() bool {
	return n.next != nil && n.next.next == nil
}
