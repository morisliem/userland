package helper

import (
	"math/rand"
	"time"

	"github.com/stretchr/testify/mock"
)

type GenerateRandomNumberMock struct {
	mock.Mock
}

func GenerateRandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999

	randomNumber := rand.Intn(max-min+1) + min

	return randomNumber
}

func (m GenerateRandomNumberMock) GenerateRandomNumber() int {
	args := m.Called()
	return args.Int(0)
}
