package bmt

// Bmt using generics allows for different sum types (like string or stdlib's sha3 [32]byte etc.)
type Bmt[T any] struct {
	sums []T
	next *Bmt[T]
	comb func(l, r T) T
}

func New[T any](combF func(l, r T) T) *Bmt[T] {
	return &Bmt[T]{
		comb: combF,
	}
}

func (b *Bmt[T]) RootSum() T {
	for {
		if b.next == nil {
			return b.sums[0]
		}
		b = b.next
	}
}

func (b *Bmt[T]) Len() int {
	return len(b.sums)
}

func (b *Bmt[T]) Add(sum T) T {
	return b.add(sum, false)
}

func (b *Bmt[T]) Size() int {
	if b.next == nil {
		return 1
	}

	return len(b.sums) + b.next.Size()
}

func (b *Bmt[T]) add(sum T, replace bool) (digest T) {
	if replace {
		b.sums = b.sums[:len(b.sums)-1]
		replace = false
	}

	b.sums = append(b.sums, sum)

	if len(b.sums) == 1 {
		return b.sums[0]
	}

	var left, right T

	length := len(b.sums)

	if length%2 == 0 {
		left = b.sums[length-2]
		right = b.sums[length-1]
		replace = true
	} else {
		left = b.sums[length-1]
		var zero T
		right = zero
	}

	combined := b.comb(left, right)

	if b.next == nil {
		b.next = &Bmt[T]{comb: b.comb}
		replace = false
	}

	return b.next.add(combined, replace)
}

func (b *Bmt[T]) Proof(i int) []T {
	return b.proof(i, nil)
}

func (b *Bmt[T]) proof(index int, acc []T) []T {
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

	if proofIndex < len(b.sums) {
		acc = append(acc, b.sums[proofIndex])
	} else {
		var zero T
		acc = append(acc, zero)
	}

	if b.isLast() {
		return acc
	}

	if left {
		index++
	}

	index /= 2

	return b.next.proof(index, acc)
}

func (b *Bmt[T]) isLast() bool {
	return b.next != nil && b.next.next == nil
}
