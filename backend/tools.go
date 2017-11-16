package main

import "math/rand"

const letters = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func createRoomNumber() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
