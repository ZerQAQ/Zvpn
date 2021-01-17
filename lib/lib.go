package lib

import (
	"crypto/rand"
	"log"
)

func ErrHandler(msg string, err error) bool {
	if err != nil {
		log.Printf("%v: %v\n", msg, err)
		return true
	} else {
		return false
	}
}

func RandomBytes(length int) []byte {
	ret := make([]byte, length)
	rand.Read(ret)
	return ret
}
