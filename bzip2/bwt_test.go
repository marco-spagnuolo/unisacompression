package bzip2

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/marco-spagnuolo/unisacompression/bzip2/bwt"
)

func TestEncode1(t *testing.T) {
	// Generate random input of length 12
	input := []byte(generateRandomString(12))

	// Call the Encode1 function
	output, err := bwt.Encode1(input)
	if err != nil {
		t.Error(err)
	}

	// Check that the output is not empty
	if len(output) == 0 {
		t.Error("Expected non-empty output")
	}
}

func TestDecode1(t *testing.T) {
	// Generate random input of length 12
	input := []byte(generateRandomString(12))

	// Encode the input using Encode1
	encoded, err := bwt.Encode1(input)
	if err != nil {
		t.Error(err)
	}

	// Call the Decode1 function
	output, err := bwt.Decode1(encoded)
	if err != nil {
		t.Error(err)
	}

	// Check that the output is as expected
	if !bytes.Equal(output, input) {
		t.Errorf("Expected output %v but got %v", input, output)
	}
}

// Helper function to generate a random string of a given length
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
