package desultory

import (
	"math/rand"
	"time"
)

var initializedRandom = false
var usedTokens map[string]bool
var randomRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func GetRandomString(length int) string {
	if !initializedRandom {
		rand.Seed(time.Now().UnixNano())
		initializedRandom = true
	}
	b := make([]rune, length)
	for i := range b {
		b[i] = randomRunes[rand.Intn(len(randomRunes))]
	}
	return string(b)
}

func GetUniqueString(length int) string {
	if usedTokens == nil {
		usedTokens = make(map[string]bool, 0)
	}
	t := GetRandomString(length)
	_, ok := usedTokens[t]
	for ok {
		t = GetRandomString(length)
		_, ok = usedTokens[t]
	}
	usedTokens[t] = true
	return t
}
