package util

import (
	"math/rand"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandomString() string {
	rand.Seed(time.Now().Unix())
	return RandomStringWithLength(RandomIntegerRange(3, 16))
}

func RandomStringWithLength(length int) string {
	b := make([]byte, length)
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandomInteger() int {
	rand.Seed(time.Now().Unix())
	return rand.Int()
}

func RandomIntegerRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func RandomFloat() float64 {
	rand.Seed(time.Now().Unix())
	return rand.Float64()
}

func RandomFloatRange(min, max float64) float64 {
	rand.Seed(time.Now().Unix())
	return min + rand.Float64()*(max-min)
}

func RandomBoolean() bool {
	cache := src.Int63()
	return cache&0x01 == 1
}
