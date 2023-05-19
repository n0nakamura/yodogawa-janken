package main

import (
	"fmt"
	"log"
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

	var playerHand string
	fmt.Print("Input: ")
	_, err := fmt.Scanf("%s", &playerHand)
	if err != nil || (playerHand != "r" && playerHand != "s" && playerHand != "p") {
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

func biasJanken() string {
	randNum := rand.Intn(100)
	if randNum < pRock {
		return "r"
	} else if randNum < (pRock + pScissors) {
		return "s"
	} else if randNum < (pRock + pScissors + pPaper) {
		return "p"
	} else {
		return "i"
	}
}

func getHandName(hand string) string {
	handNames := map[string]string{
		"r": "Rock",
		"s": "Scissors",
		"p": "Paper",
		"i": "Invincible",
	}

	name, ok := handNames[hand]
	if !ok {
		log.Fatal("Unknown hand")
	}

	return name
}

func judge(playerHand, yodogawaHand string) int {
	if playerHand == yodogawaHand {
		return 0
	}
	if (playerHand == "r" && yodogawaHand == "s") ||
		(playerHand == "s" && yodogawaHand == "p") ||
		(playerHand == "p" && yodogawaHand == "r") {
		return 1
	}
	return -1
}
