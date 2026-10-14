package base62

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base = int64(len(charset))

// Encode converts an integer ID to a Base62 string
func Encode(num int64) string {
	if num == 0 {
		return string(charset[0])
	}

	var sb strings.Builder
	for num > 0 {
		rem := num % base
		sb.WriteByte(charset[rem])
		num = num / base
	}

	// Reverse the string
	bytes := []byte(sb.String())
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return string(bytes)
}

// GenerateRandom generates a random alphanumeric string of specified length
func GenerateRandom(length int) (string, error) {
	bytes := make([]byte, length)
	charsetLen := big.NewInt(base)

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		bytes[i] = charset[randomIndex.Int64()]
	}

	return string(bytes), nil
}
