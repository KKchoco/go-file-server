package main

import "math/rand"
import "time"

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func RandString(length int) string {
	letters := []byte("ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz1234567890-_~")

	// Make a slice of bytes, allocating n bytes
	bytes := make([]byte, length)

	// Iterate through the slice, setting each byte to a random letter
	for i := range bytes {
		rand.Seed(time.Now().UnixNano())
		bytes[i] = letters[rand.Intn(len(letters))]
	}

	return string(bytes)
}
