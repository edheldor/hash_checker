package main

import (
	"crypto/sha256"
	"fmt"
	//"os"
	//"path/filepath"
)

func main() {
    smaple := file{name: "dsd", path: "dsfdsf"}
    smaple.calclHash()
    fmt.Println(smaple.hash)

}


type file struct {
	name string
	path string
	hash string
}

 func (item *file)  calclHash() {
	hash := sha256.New()
	hash.Write([]byte("HELLO"))
	item.hash = fmt.Sprintf("%x", hash.Sum(nil))


}