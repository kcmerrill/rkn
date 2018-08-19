package rkn

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// DB ...
type DB struct {
	// private
	toDisk     chan string          // does this need to write to disk?
	sync       chan bool            // when to flush to disk
	close      chan bool            // close the db cleanly
	exit       chan bool            // exit the db properly
	docs       map[string]*Version  // markers to location of the documents in the file
	cache      map[string]*Document // document version cache
	segment    string
	fileWriter *os.File            // handle to write to the file
	fileMarker int64               // marks our currently location in the file
	keyScrub   func(string) string // key validation
	sync.Mutex                     // global lock
}

// Open ...
func Open(store string, d *DB) (*DB, error) {
	// init our maps
	d.docs = make(map[string]*Version)
	d.cache = make(map[string]*Document)

	// init our channels
	d.toDisk = make(chan string)
	d.close = make(chan bool)
	d.exit = make(chan bool)
	d.sync = make(chan bool)

	// create a valid key string replacer
	d.keyScrub = strings.NewReplacer("{", "_", "}", "_", ":", "_", "|", "_").Replace

	// open up our file for writing ...
	var err error
	d.fileWriter, err = os.OpenFile(store, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return d, err
	}

	// start our background processor
	go d.BGProcesser()

	// return the goods
	return d, nil
}

// Close ...
func (d *DB) Close() {
	d.close <- true
	<-d.exit
}

// Save ...
func (d *DB) Save(doc *Document) {
	key, ts, contents := doc.toString()
	d.Lock()
	d.cache[key] = doc
	d.Unlock()
	d.toDisk <- fmt.Sprintf("{%s|%s|%d}%s", d.keyScrub(key), ts, len(contents), contents)
}

// Fetch ...
func (d *DB) Fetch(key string) (doc *Document, exists bool) {
	return d.fetch(key, 0)
}

// FetchVersion ...
func (d *DB) FetchVersion(key string, version int) (doc *Document, exists bool) {
	return d.fetch(key, version)
}

// internal fetch, w/version
func (d *DB) fetch(key string, version int) (doc *Document, exists bool) {
	key = d.keyScrub(key)
	d.Lock()
	doc, exists = d.cache[key]
	d.Unlock()

	if !exists {
		// bleh, not in the cache ...
		/*d.Lock()
		v, exists := d.docs[key]
		d.Unlock()
		if exists {
			doc := d.retrieveDoc(v.Get(version))
		}
		*/
	}
	return doc, exists
}

// Doc returns a newly created document
func (d *DB) Doc(key, contents string) *Document {
	return &Document{
		Key:       d.keyScrub(key),
		Contents:  contents,
		TimeStamp: time.Now(),
	}
}

// BGProcesser is our back ground processor
func (d *DB) BGProcesser() {
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				// we want to sync every second
				d.sync <- true
			}
		}
	}()

	for {
		select {
		case <-d.sync:
			d.fileWriter.Sync()
		case <-d.close:
			fmt.Println("close")
			d.fileWriter.Sync()
			d.exit <- true
		case doc := <-d.toDisk:
			d.fileWriter.WriteString(doc + "\n")
		}
	}
}
