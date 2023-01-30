package bzip2

import (
	"bytes"
	"sort"
)

// The Burrows-Wheeler Transform implementation used here is based on the
// Suffix Array by Induced Sorting (SA-IS) 
// 
//
// The SA-IS algorithm runs in O(n) and outputs a Suffix Array. There is a
// mathematical relationship between Suffix Arrays and the Burrows-Wheeler
// Transform, such that a SA can be converted to a BWT in O(n) time.
//
// References:
//	https://github.com/cscott/compressjs/blob/master/lib/BWT.js
//	
type burrowsWheelerTransform struct {
	buf  []byte
	sa   []int
	perm []uint32
}

// func (bwt *burrowsWheelerTransform) Encode(buf []byte) (ptr int) {
//  	if len(buf) == 0 {
//  		return -1
//  	}



// 	// Step 1: Concatenate the input string to itself so that we can use the
//  	// suffix array algorithm for bzip2's variant of BWT.
//  	n := len(buf)
//  	bwt.buf = append(append(bwt.buf[:0], buf...), buf...)
//  	if cap(bwt.sa) < 2*n {
//  		bwt.sa = make([]int, 2*n)
//  	}
//  	t := bwt.buf[:2*n]
//  	sa := bwt.sa[:2*n]

//  	// Step 2: Compute the suffix array (SA). The input string, t, will not be
//  	// modified, while the results will be written to the output, sa.
//  	sais.ComputeSA(t, sa)

//  	// Step 3: Convert the SA to a BWT. Since ComputeSA does not mutate the
//  	// input, we have two copies of the input; in buf and buf2. Thus, we write
//  	// the transformation to buf, while using buf2.
//  	var j int
//  	buf2 := t[n:]
//  	for _, i := range sa {
//  		if i < n {
//  			if i == 0 {
//  				ptr = j
//  				i = n
//  			}
//  			buf[j] = buf2[i-1]
//  			j++
//  		}
//  	}
//  	return ptr
//  }

// func (bwt *burrowsWheelerTransform) Decode(buf []byte, ptr int) {
// 	if len(buf) == 0 {
// 		return
// 	}

// 	// Step 1: Compute cumm, where cumm[ch] reports the total number of
// 	// characters that precede the character ch in the alphabet.
// 	var cumm [256]int
// 	for _, v := range buf {
// 		cumm[v]++
// 	}
// 	var sum int
// 	for i, v := range cumm {
// 		cumm[i] = sum
// 		sum += v
// 	}

// 	// Step 2: Compute perm, where perm[ptr] contains a pointer to the next
// 	// byte in buf and the next pointer in perm itself.
// 	if cap(bwt.perm) < len(buf) {
// 		bwt.perm = make([]uint32, len(buf))
// 	}
// 	perm := bwt.perm[:len(buf)]
// 	for i, b := range buf {
// 		perm[cumm[b]] = uint32(i)
// 		cumm[b]++
// 	}

// 	// Step 3: Follow each pointer in perm to the next byte, starting with the
// 	// origin pointer.
// 	if cap(bwt.buf) < len(buf) {
// 		bwt.buf = make([]byte, len(buf))
// 	}
// 	buf2 := bwt.buf[:len(buf)]
// 	i := perm[ptr]
// 	for j := range buf2 {
// 		buf2[j] = buf[i]
// 		i = perm[i]
// 	}
// 	copy(buf, buf2)
// }

func (bwt *burrowsWheelerTransform) Encode1(buf []byte) (ptr int) {
	if len(buf) == 0 {
		return -1
	}

	// Create a list of rotations of the src slice
	rotations := make([][]byte, len(buf))
	for i := range buf {
		rotation := make([]byte, len(buf))
		copy(rotation[:i], buf[len(buf)-i:])
		copy(rotation[i:], buf[:len(buf)-i])
		rotations[i] = rotation
	}

	// Sort the rotations lexicographically
	sort.Slice(rotations, func(i, j int) bool {
		return bytes.Compare(rotations[i], rotations[j]) == -1
	})

	// Write the last characters of each rotation to the dst slice
	bwt.perm = make([]uint32, len(buf))
	for i, rotation := range rotations {
		bwt.perm[i] = uint32(rotation[len(rotation)-1])
		if bytes.Equal(rotation, buf) {
			ptr = i
		}
	}

	return ptr
}

func (bwt *burrowsWheelerTransform) Decode1(buf []byte, ptr int) {
	// Compute the inverse permutation of the BWT
	invPerm := make([]int, len(bwt.perm))
	for i, p := range bwt.perm {
		invPerm[p] = i
	}

	// Reconstruct the original string using the inverse permutation
	for i, p := range invPerm {
		if p == ptr {
			buf[i] = 0
		} else {
			buf[i] = byte(bwt.perm[p])
		}
	}
}
