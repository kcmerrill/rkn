package rkn

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/snappy"
)

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
	d.TimeStamp = time.Now()
	d.Unlock()

	return d.TimeStamp
}

// toString()
func (d *Document) toString() (key string, timestamp string, compressed string) {
	d.Lock()
	key = strings.Replace(d.Key, "|", "-", -1)
	timestamp = strconv.FormatInt(d.TimeStamp.Unix(), 10)
	compressed = string(snappy.Encode(nil, []byte(d.Contents)))
	d.Unlock()

	return key, timestamp, compressed
}
