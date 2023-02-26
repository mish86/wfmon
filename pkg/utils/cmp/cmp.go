package cmp

import (
	"wfmon/pkg/utils/conv"

	"golang.org/x/exp/constraints"
)

func Compare[T constraints.Ordered](a, b T) int {
	return conv.BoolToInt(a > b) - conv.BoolToInt(a < b)
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Nvl[T any](expr bool, t, f T) T {
	if expr {
		return t
	}
	return f
}
