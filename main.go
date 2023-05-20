package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Outcome uint

const (
	WIN Outcome = iota
	LOSE
	DRAW
)

var outcomeNameMap = map[Outcome]string{
	WIN:  "YOU WIN",
	LOSE: "YOU LOSE",
	DRAW: "DRAW",
}

func (o Outcome) String() string {
	if s, ok := outcomeNameMap[o]; ok {
		return s
	}
	return "Unknown status"
}

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
	case DRAW:
		fmt.Println(DRAW)
	case WIN:
		fmt.Println(WIN)
	case LOSE:
		fmt.Println(LOSE)
	}
}

func biasJanken() string {
	var handProbabilities = map[string]int{
		"r": 26,
		"s": 32,
		"p": 37,
		"i": 5,
	}

	var allProb = 0
	for _, probability := range handProbabilities {
		allProb += probability
	}
	randNum := rand.Intn(allProb)

	var cumulativeProb = 0
	for hand, probability := range handProbabilities {
		cumulativeProb += probability
		if randNum < cumulativeProb {
			return hand
		}
	}
	return ""
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

func judge(playerHand, yodogawaHand string) Outcome {
	if playerHand == yodogawaHand {
		return DRAW
	}
	if (playerHand == "r" && yodogawaHand == "s") ||
		(playerHand == "s" && yodogawaHand == "p") ||
		(playerHand == "p" && yodogawaHand == "r") {
		return WIN
	}
	return LOSE
}
