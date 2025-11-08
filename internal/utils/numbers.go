package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateRandomNumber(length int) int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%d", r.Intn(9)+1))
	for i := 1; i < length; i++ {
		sb.WriteString(fmt.Sprintf("%d", r.Intn(10)))
	}

	num, _ := strconv.Atoi(sb.String())

	return num
}
