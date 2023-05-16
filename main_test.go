package yodogawajanken

import (
	"fmt"
	"testing"
)

func TestBiasJanken(t *testing.T) {
	const c = 100000
	r, s, p, i := 0.0, 0.0, 0.0, 0.0
	for j := 0; j < c; j++ {
		hand := biasJanken()
		switch hand {
		case 0:
			r++
		case 1:
			s++
		case 2:
			p++
		case 3:
			i++
		default:
			t.Error("Unknown hand")
		}
	}
	fmt.Println("r=", r/c*100, "% s=", s/c*100, "% p=", p/c*100, "% i=", i/c*100, "%")

	if !((i < r) && (r < s) && (s < p)) {
		t.Error("Inappropriate rate of occurrence")
	}
}
