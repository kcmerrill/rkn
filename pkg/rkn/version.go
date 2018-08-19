package rkn

import (
	"sync"
)

// Version will hold the location to each version in the db
type Version struct {
	sync.Mutex
	versions []int64
}

// Get ...
func (v *Version) Get(version int) int64 {
	if version > 0 && version <= len(v.versions)-1 {
		// version exists ... lets get it!
		v.Lock()
		val := v.versions[version]
		v.Unlock()
		return val
	}
	return v.Latest()
}

// Add a version
func (v *Version) Add(version int64) {
	v.versions = append(v.versions, version)
}

// Latest will return the marker for the latest version
func (v *Version) Latest() (version int64) {
	// return the latest and greatest
	v.Lock()
	version = v.versions[len(v.versions)-1]
	v.Unlock()
	return version
}
