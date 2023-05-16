package yodogawajanken

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	pRock     = 26
	pScissors = 32
	pPaper    = 37
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var playerHand int
	fmt.Print("Input: ")
	_, err := fmt.Scanf("%d", &playerHand)
	if err != nil || playerHand < 0 || playerHand > 2 {
		fmt.Println("Error")
		return
	}

	yodogawaHand := biasJanken()

	fmt.Println("Your hand: ", getHandName(playerHand))
	fmt.Println("Yodogawa-san hand: ", getHandName(yodogawaHand))

	result := judge(playerHand, yodogawaHand)
	switch result {
	case 0:
		fmt.Println("DRAW")
	case 1:
		fmt.Println("YOU WIN")
	case -1:
		fmt.Println("YOU LOSE")
	}
}

func biasJanken() int {
	randNum := rand.Intn(100)
	if randNum < pRock {
		return 0
	} else if randNum < (pRock + pScissors) {
		return 1
	} else if randNum < (pRock + pScissors + pPaper) {
		return 2
	} else {
		return 3
	}
}

func getHandName(hand int) string {
	switch hand {
	case 0:
		return "Rock"
	case 1:
		return "Scissors"
	case 2:
		return "Paper"
	case 3:
		return "Invincible"
	default:
		return "Unknown hand"
	}
}

func judge(playerHand, yodogawaHand int) int {
	if playerHand == yodogawaHand {
		return 0
	}
	if (playerHand == 0 && yodogawaHand == 1) ||
		(playerHand == 1 && yodogawaHand == 2) ||
		(playerHand == 2 && yodogawaHand == 0) {
		return 1
	}
	return -1
}

