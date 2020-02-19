// Copyright (c) 2012, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package filter provides interface and implementation of probabilistic
// data structure.
//
// The filter is resposible for creating small filter from a set of keys.
// These filter will then used to test whether a key is a member of the set.
// In ptcy cases, a filter can cut down the number of disk seeks from a
// handful to a single disk seek per DB.Get call.
package filter

// Buffer is the interface that wraps basic Alloc, Write and WriteByte methods.
type Buffer interface {
	// Alloc allocs n bytes of slice from the buffer. This also advancing
	// write offset.
	Alloc(n int) []byte

	// Write appends the contents of p to the buffer.
	Write(p []byte) (n int, err error)

	// WriteByte appends the byte c to the buffer.
	WriteByte(c byte) error
}

// Filter is the filter.
type Filter interface {
	// Name returns the name of this policy.
	//
	// Note that if the filter encoding changes in an incompatible way,
	// the name returned by this method must be changed. Otherwise, old
	// incompatible filters may be passed to methods of this type.
	Name() string

	// NewGenerator creates a new filter generator.
	NewGenerator() FilterGenerator

	// Contains returns true if the filter contains the given key.
	//
	// The filter are filters generated by the filter generator.
	Contains(filter, key []byte) bool
}

// FilterGenerator is the filter generator.
type FilterGenerator interface {
	// Add adds a key to the filter generator.
	//
	// The key may become invalid after call to this method end, therefor
	// key must be copied if implementation require keeping key for later
	// use. The key should not modified directly, doing so may cause
	// undefined results.
	Add(key []byte)

	// Generate generates filters based on keys passed so far. After call
	// to Generate the filter generator maybe resetted, depends on implementation.
	Generate(b Buffer)
}
