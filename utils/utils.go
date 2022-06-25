package utils

import (
	"math/rand"
	"os"
	"time"
)

func GetValueFromEnv(name string, defaultVal string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return defaultVal
}

var source = rand.NewSource(time.Now().UnixNano())

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[source.Int63()%int64(len(charset))]
	}
	return string(b)
}
