package main

import "math/rand"

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func RandString(length int) string {
	letters := []byte("ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz")

	// Make a slice of bytes, allocating n bytes
	bytes := make([]byte, length)

	// Iterate through the slice, setting each byte to a random letter
	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}

	return string(bytes)
}
