package rkn

import (
	"fmt"
	"sync"
)

// DB ...
type DB struct {
	// private
	toDisk   chan string // does this need to write to disk?
	flush    chan bool   // when to flush to disk
	inMemory bool        // does this persist to memory?
	cache    map[string]*Document

	// global db lock
	sync.Mutex
}

// Open ...
func Open(store string) (*DB, error) {
	d := &DB{}

	// are we in memory?
	d.inMemory = d.isInMemory(store)

	// init our maps
	d.segments = make(map[string]int64)

	// init our channels
	d.toDisk = make(chan string)

	go d.Broker()
	return d, nil
}

// Save ...
func (d *DB) Save(doc *Document) {
	key, ts, contents := doc.toString()
	d.toDisk <- fmt.Sprintf("{%s|%s|%d}%s", key, ts, len(contents), contents)
}

// Fetch ...
func (d *DB) Fetch(key string) *Document {
	return &Document{}
}

// isInMemory() determines if the store is in memory
func (d *DB) isInMemory(store string) bool {
	return store == ":in-memory:"
}

// Broker ...
func (d *DB) Broker() {
	for {
		select {}
	}
}
