package main

import (
	"fmt"
	"regexp"
	"testing"
)

func TestJanken(t *testing.T) {
	const c = 100000
	r, s, p, i := 0.0, 0.0, 0.0, 0.0

	re := regexp.MustCompile(`Yodogawa-san hand: . (.+)`)

	for j := 0; j < c; j++ {
		hand, _ := janken("R")
		match := re.FindStringSubmatch(hand)
		switch match[1] {
		case "Rock":
			r++
		case "Scissors":
			s++
		case "Paper":
			p++
		case "Invincible":
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

// func TestGetPlayerHand(t *testing.T) {
// 	for hand, patterns := range handPatterns {
// 		for _, h := range patterns {
// 			if result, err := getPlayerHand(string(h)); (hand != result) || (err != nil) {
// 				t.Error("Invalid return value of the function getPlayerHand within the expected pattern range.")
// 			}
// 		}
// 	}

// 	if _, err := getPlayerHand("あ"); err == nil {
// 		t.Error("No errors have occurred for unexpected patterns.")
// 	}
// }
