// Copyright (c) 2014 Josh Rickmar.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Package bitset provides bitset implementations for bit packing binary
// values into words and bytes.
//
// A bitset, while logically equivalent to a []bool, is often preferable
// over a []bool due to the space and time efficiency of bit packing binary
// values.  They are typically more space efficient than a []bool since
// (although implementation specifc) bools are typically machine word size.
// On 64-bit architectures, this can result in a bitset using about 64 times
// less memory to store the same boolean values as a []bool.  While bitsets
// introduce bitshifting overhead for gets and sets unnecessary for a []bool,
// they are usually still more performant than a []bool due to the smaller data
// structure  being more cache friendly.
//
// This package contains three bitset implementations: Words for efficiency,
// Bytes for situations where bitsets must be serialized or deserialized,
// and Spare for when memory efficiency is the most important factor when
// working with sparse datasets.
package bitset
