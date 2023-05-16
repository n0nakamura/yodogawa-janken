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
		case "r":
			r++
		case "s":
			s++
		case "p":
			p++
		case "i":
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
