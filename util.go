package sfapi

import (
	"fmt"
	"log"
	"math"
	"strings"
)

func Min(x, y int64) int64 {
	if x < y {
		return x
	}

	return y
}

func Dollars(cents int64) string {
	d := float64(cents) / 100.00

	return fmt.Sprintf("$%.2f", d)
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}

	return x
}

func PadString(s string, width int) (result string) {
	if len(s) < width {
		result = s + strings.Repeat(" ", width-len(s))
	}

	return
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func FloatEquals(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
