package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	go func() {
		f, _ := os.OpenFile("/tmp/dat.txt", os.O_CREATE|os.O_RDONLY, 0644)
		for {
			b := make([]byte, 1)
			f.ReadAt(b, int64(1))
			fmt.Println(string(b))
		}
	}()
	go func() {
		f, _ := os.OpenFile("/tmp/dat.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		for {
			f.WriteString("Here is a string, here is another string\n")
		}
	}()
	<-time.After(10 * time.Second)
	fmt.Println("Finished!")
}
