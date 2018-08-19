package main

import (
	"time"

	"github.com/kcmerrill/rkn/pkg/rkn"
)

func main() {
	db, _ := rkn.Open("/tmp/something.txt", &rkn.DB{})
	defer db.Close()

	doc := db.Doc("key", "bingo was his nameo was was was")
	go func() {
		for {
			db.Save(doc)
		}
	}()
	<-time.After(time.Second)
}
