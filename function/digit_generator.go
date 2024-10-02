package function

import (
	"fmt"
	"math/rand"
	"time"
)

func DigitGenerator() string {
	var result string

	rand.Seed(time.Now().UnixNano())

	for i := 0;i < 6;i++{
		result += fmt.Sprintf("%d", rand.Intn(10))
	}

	return result
}