package utils

import (
	"fmt"
	"math/rand"
)

func RandomColor() string {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)

	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

func FlattenMap[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
