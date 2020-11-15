package util

import (
	"math/rand"
	"net"
	"strconv"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateRandomString -> Ron Seal
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// IsAnInteger I mean, if I don't know what this does then lol.
func IsAnInteger(input string) bool {
	_, err := strconv.Atoi(input)
	if err != nil {
		return false
	}
	return true
}

// IsLegalIPAddress is used when a user tries setting an LHOST value. Quick sanity check incase of derps
func IsLegalIPAddress(input string) bool {
	parsedIP := net.ParseIP(input)
	if parsedIP.To4() == nil {
		return false
	}
	return true
}

// IsLegalPortNumber is used when a user tries setting an LPORT value. Quick sanity check incase of derps
func IsLegalPortNumber(input string) bool {
	val, err := strconv.Atoi(input)
	if err != nil {
		return false
	}

	if val < 0 || val > 65535 { // right? computers are hard
		return false
	}
	return true
}

// Xor function used to xor the shellcode before adding it to the implant source
func Xor(data string, key string) string {
	output := ""

	for i := 0; i < len(data); i++ {
		output += string(data[i] ^ key[i%len(key)])
	}

	return output
}
