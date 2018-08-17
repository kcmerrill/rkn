package rkn

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/snappy"
)

// Doc returns a newly created document
func Doc(key, contents string) *Document {
	return &Document{
		Key:       key,
		Contents:  contents,
		TimeStamp: time.Now(),
	}
}

// Document ...
type Document struct {
	TimeStamp time.Time `json:"timestamp"`
	Contents  string    `json:"contents"`
	Key       string    `json:"key"`
	sync.Mutex
}

// Touch ...
func (d *Document) Touch() time.Time {
	d.Lock()
	defer d.Unlock()

	d.TimeStamp = time.Now()
	return d.TimeStamp
}

// toString()
func (d *Document) toString() (string, string, string) {
	d.Lock()
	defer d.Unlock()

	compressed := snappy.Encode(nil, []byte(d.Contents))
	return strings.Replace(d.Key, "|", "-", -1), strconv.FormatInt(d.TimeStamp.Unix(), 10), string(compressed)
}
