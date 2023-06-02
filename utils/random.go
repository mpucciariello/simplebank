package utils

import (
	"math/rand"
	"strings"
)

// generates random values ir order to avoid register duplication while testing

const alphabet = "qwertyuiopasdfhjklzxcvbnm"

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(n int) string {
	var sb strings.Builder
	l := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(l)]
		sb.WriteByte(c)
	}

	return strings.Title(sb.String())
}

func RandomOwner() string {
	return randomString(6)
}

func RandomBalance() int64 {
	return randomInt(0, 10000)
}

func RandomCurrency() string {
	currencies := []string{"ARS", "USD", "EUR"}
	return currencies[rand.Intn(len(currencies))]
}
