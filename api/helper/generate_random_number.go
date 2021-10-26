package helper

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999

	randomNumber := rand.Intn(max-min+1) + min

	return randomNumber
}
