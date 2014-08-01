// Copyright (c) 2014 Josh Rickmar.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package bitset

const (
	// wordBits is the total number of bits that make up a word.
	wordBits = 32 << uint(^uintptr(0)>>63)

	// wordModMask is the maximum number of indexes a 1 may be left
	// shifted by before the value overflows a word.  It is equal to one
	// less than the number of bits per word.
	//
	// This package uses this value to calculate bit indexes within a
	// word as it is quite a bit more efficient to perform a bitwise AND
	// with this rather than using the modulus operator (n&31 == n%32,
	// and n%63 == n%64).
	wordModMask = wordBits - 1

	// wordShift is the number of bits to perform a right shift of a bit
	// index by to get the word index of a bit in the bitset.  It is
	// functionally equivalent to integer dividing the bit index by
	// wordBits, but is a bit more efficient to calculate.  The value is
	// equal to the log2(wordBits), or the single bit index set in a
	// word to create the value wordBits.  On machines where a word is
	// 32-bits, this value is 5.  On machines with 64-bit sized words,
	// this value is 6.  It is calculated as follows:
	//
	//      On 32-bit machines (results in 5):
	//        Given the word size 32:  0b00100000
	//        Add the value 128:       0b10100000
	//        And right shift by 5:    0b00000101
	//
	//      On 64-bit machines (results in 6):
	//        Given the word size 64:  0b01000000
	//        Add the value 128:       0b11000000
	//        And right shift by 5:    0b00000110
	wordShift = (1<<7 + wordBits) >> 5

	// byteModMask is the maximum number of indexes a 1 may be left
	// shifted by before the value overflows a byte.  It is equal to one
	// less than the number of bits per byte.
	//
	// This package uses this value to calcualte all bit indexes within
	// a byte as it is quite a bit more efficient to perform a bitwise
	// AND with this rather than using the modulus operator (n&7 == n%8).
	byteModMask = 7 // 0b0000111

	// byteShift is the number of bits to perform a right shift of a bit
	// index to get the byte index in a bitset.  It is functionally
	// equivalent to integer dividing by 8 bits per byte, but is a bit
	// more efficient to calculate.
	byteShift = 3
)

// BitSet defines the method set of a bitset.  Bitsets allow for bit
// packing of binary values to words or bytes for space and time
// efficiency.
//
// The Grow methods of Words and Bytes are not part of this interface.
type BitSet interface {
	Get(i uint) bool
	Set(i uint)
	Unset(i uint)
	SetBool(i uint, b bool)
}

// Words represents a bitset backed by a word slice.  Words bitsets are
// designed for effeciency and do not automatically grow for indexed values
// outside of the allocated range.  The Grow method is provided if it is
// necessary to grow a Words bitset beyond its initial allocation.
//
// The len of a Words is the number of words in the set.  Multiplying by
// the machine word size will result in the number of bits the set can hold.
type Words []uintptr

// NewWords returns a new bitset that is capable of holding numBits number
// of binary values.  All words in the bitset are zeroed and each bit is
// therefore considered unset.
func NewWords(numBits uint) Words {
	return make(Words, (numBits+wordModMask)>>wordShift)
}

// Get returns whether the bit at index i is set or not.  This method will
// panic if the index results in a word index that exceeds the number of
// words held by the bitset.
func (w Words) Get(i uint) bool {
	return w[i>>wordShift]&(1<<(i&wordModMask)) != 0
}

// Set sets the bit at index i.  This method will panic if the index results
// in a word index that exceeds the number of words held by the bitset.
func (w Words) Set(i uint) {
	w[i>>wordShift] |= 1 << (i & wordModMask)
}

// Unset unsets the bit at index i.  This method will panic if the index
// results in a word index that exceeds the number of words held by the
// bitset.
func (w Words) Unset(i uint) {
	w[i>>wordShift] &^= 1 << (i & wordModMask)
}

// SetBool sets or unsets the bit at index i depending on the value of b.
// This method will panic if the index results in a word index that exceeds
// the number of words held by the bitset.
func (w Words) SetBool(i uint, b bool) {
	if b {
		w.Set(i)
		return
	}
	w.Unset(i)
}

// Grow ensures that the bitset w is large enough to hold numBits number of
// bits, potentially appending to and/or reallocating the slice if the
// current length is not sufficient.
func (w *Words) Grow(numBits uint) {
	words := *w
	targetLen := (numBits + wordModMask) >> wordShift
	if missing := targetLen - uint(len(words)); missing != 0 {
		*w = append(words, make(Words, missing)...)
	}
}

// Bytes represents a bitset backed by a bytes slice.  Bytes bitsets,
// while designed for efficiency, are slightly less effecient to use
// than Words bitsets, since word-sized data is faster to manipulate.
// However, Bytes have the nice property of easily and portably being
// (de)serialized to or from an io.Reader or io.Writer.  Like a Words,
// Bytes bitsets do not automatically grow for indexed values outside
// of the allocated range.  The Grow method is provided if it is
// necessary to grow a Bytes bitset beyond its initial allocation.
//
// The len of a Bytes is the number of bytes in the set.  Multiplying by
// eight will result in the number of bits the set can hold.
type Bytes []byte

// NewBytes returns a new bitset that is capable of holding numBits number
// of binary values.  All bytes in the bitset are zeroed and each bit is
// therefore considered unset.
func NewBytes(numBits uint) Bytes {
	return make(Bytes, (numBits+byteModMask)>>byteShift)
}

// Get returns whether the bit at index i is set or not.  This method will
// panic if the index results in a byte index that exceeds the number of
// bytes held by the bitset.
func (s Bytes) Get(i uint) bool {
	return s[i>>byteShift]&(1<<(i&byteModMask)) != 0
}

// Set sets the bit at index i.  This method will panic if the index results
// in a byte index that exceeds the number of a bytes held by the bitset.
func (s Bytes) Set(i uint) {
	s[i>>byteShift] |= 1 << (i & byteModMask)
}

// Unset unsets the bit at index i.  This method will panc if the index
// results in a byte index that exceeds the number of bytes held by the
// bitset.
func (s Bytes) Unset(i uint) {
	s[i>>byteShift] &^= 1 << (i & byteModMask)
}

// SetBool sets or unsets the bit at index i depending on the value of b.
// This method will panic if the index results in a byte index that exceeds
// the nubmer of bytes held by the bitset.
func (s Bytes) SetBool(i uint, b bool) {
	if b {
		s.Set(i)
		return
	}
	s.Unset(i)
}

// Grow ensures that the bitset s is large enough to hold numBits number of
// bits, potentially appending to and/or reallocating the slice if the
// current length is not sufficient.
func (s *Bytes) Grow(numBits uint) {
	bytes := *s
	targetLen := (numBits + byteModMask) >> byteShift
	if missing := targetLen - uint(len(bytes)); missing != 0 {
		*s = append(bytes, make(Bytes, missing)...)
	}
}

// Sparse is a memory effecient bitset for sparsly-distributed set bits.
// Unlike a Words or Bytes which requires each word or byte between 0 and
// the highest index to be allocated, a Sparse only holds the words which
// contain set bits.  Additionally, Spare is the only BitSet implementation
// from this package which will dynamically expand and shrink as bits are
// set and unset.
//
// As the map is unordered and there is no obvious way to (de)serialize this
// structure, no byte implementation is provided, and all map values are
// machine word sized.
//
// As Sparse bitsets are backed by a map, getting and setting bits are
// orders of magnitude slower than other slice-backed bitsets and should
// only be used with sparse datasets and when memory effeciency is a
// top concern.  It is highly recommended to benchmark this type against
// the other bitsets using realistic sample data before using this type
// in an application.
//
// New Spare bitsets can be created using the builtin make function.
type Sparse map[uint]uintptr

// Get returns whether the bit at index i is set or not.
func (s Sparse) Get(i uint) bool {
	return s[i>>wordShift]&(1<<(i&wordModMask)) != 0
}

// Set sets the bit at index i.  A word insert is performed if if no bits
// of this word have been previously set.
func (s Sparse) Set(i uint) {
	s[i>>wordShift] |= 1 << (i & wordModMask)
}

// Unset unsets the bit at index i.  If all bits for a given word have are
// unset, the word is removed from the set, and future calls to Get will
// return false for all bits from this word.
func (s Sparse) Unset(i uint) {
	wordKey := i >> wordShift
	word, ok := s[wordKey]
	if !ok {
		return
	}
	word &^= 1 << (i & wordModMask)
	if word == 0 {
		delete(s, wordKey)
	} else {
		s[wordKey] = word
	}
}

// SetBool sets the bit at index i if b is true, otherwise the bit is unset.
// see the comments for the get and set methods for the memory allocation
// rules that are followed when getting or setting bits in a Sparse bitset.
func (s Sparse) SetBool(i uint, b bool) {
	if b {
		s.Set(i)
		return
	}
	s.Unset(i)
}
