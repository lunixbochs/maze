package main

import (
	"crypto/rand"
	"log"
)

func randByte() int {
	bytes := []byte{0}
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return int(bytes[0])
}
