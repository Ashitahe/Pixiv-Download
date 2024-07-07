package random

import (
	"math/rand"
	"time"
)

func GetRandomElement(list []string) string {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return list[r.Intn(len(list))]
}