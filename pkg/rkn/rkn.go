package rkn

import (
	"os"
)

// Rkn ...
type Rkn struct {
	// private
	toDiskHandle os.File
}

// Open ...
func (r *Rkn) Open(store string) {

}

// Save ...
func (r *Rkn) Save(d *Document) bool {
	return true
}
