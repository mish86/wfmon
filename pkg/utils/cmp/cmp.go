package cmp

import (
	"wfmon/pkg/utils/conv"

	"golang.org/x/exp/constraints"
)

func Compare[T constraints.Ordered](a, b T) int {
	return conv.BoolToInt(a > b) - conv.BoolToInt(a < b)
}
