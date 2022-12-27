
package bzip2

import "github.com/marco-spagnuolo/unisacompression/bzip2/internal/sais"

// The Burrows-Wheeler Transform implementation used here is based on the
// Suffix Array by Induced Sorting (SA-IS) methodology by Nong, Zhang, and Chan.
// This implementation uses the sais algorithm originally written by Yuta Mori.
//
// The SA-IS algorithm runs in O(n) and outputs a Suffix Array. There is a
// mathematical relationship between Suffix Arrays and the Burrows-Wheeler
// Transform, such that a SA can be converted to a BWT in O(n) time.
//
// References:
//	http://www.hpl.hp.com/techreports/Compaq-DEC/SRC-RR-124.pdf
//	https://github.com/cscott/compressjs/blob/master/lib/BWT.js
//	https://www.quora.com/How-can-I-optimize-burrows-wheeler-transform-and-inverse-transform-to-work-in-O-n-time-O-n-space
type burrowsWheelerTransform struct {
	buf  []byte
	sa   []int
	perm []uint32
}


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
