package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

// generates random values ir order to avoid register duplication while testing

const alphabet = "qwertyuiopasdfhjklzxcvbnm"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	l := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(l)]
		sb.WriteByte(c)
	}

	return strings.Title(sb.String())
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%v@mail.com", RandomString(8))
}

func RandomBalance() int64 {
	return RandomInt(0, 10000)
}

func RandomCurrency() string {
	currencies := []string{ARS, EUR, USD}
	return currencies[rand.Intn(len(currencies))]
}
