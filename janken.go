package main

import (
	"log"
	"math/rand"
	"time"
)

type Outcome uint

const (
	WIN Outcome = iota
	LOSE
	DRAW
	LOVE
	HANDSHAKE
)

var outcomeNameMap = map[Outcome]string{
	WIN:       "YOU WIN",
	LOSE:      "YOU LOSE",
	DRAW:      "DRAW",
	LOVE:      "BIG LOVE...🤟",
	HANDSHAKE: "SHAKING... 🤝",
}

func (o Outcome) String() string {
	if s, ok := outcomeNameMap[o]; ok {
		return s
	}
	return "Unknown status"
}

type Hand uint

const (
	ROCK Hand = iota
	SCISSORS
	PAPER
	INVINCIBLE
	HLOVE
	HHANDSHAKE
	OTHER
)

var handNames = map[Hand]string{
	ROCK:       "✊ Rock",
	SCISSORS:   "✌ Scissors",
	PAPER:      "🖐 Paper",
	INVINCIBLE: "👉 Invincible",
	HLOVE:      "🤟",
	HHANDSHAKE: "🤝",
	OTHER:      "🤔",
}

func (h Hand) String() string {
	if s, ok := handNames[h]; ok {
		return s
	}
	return "Unknown hand"
}

func getPlayerHand(playerHand string) Hand {
	switch playerHand {
	case "R", "✊", "👊", "🤛", "🤜", "💪":
		return ROCK
	case "S", "✌", "🤞":
		return SCISSORS
	case "P", "🖐", "✋", "🤚", "🖖", "🫲", "🫱", "🫳", "🫴", "👋", "👐", "🤲", "🤗":
		return PAPER
	case "🤟":
		return HLOVE
	case "🤝":
		return HHANDSHAKE
	case "👌", "🤌", "🤏", "🤘", "🤙", "👈", "👉", "👆", "👇", "☝", "👍", "👏", "🙏":
		return OTHER
	default:
		log.Fatalf("Invalid hand")
		panic("panic")
	}
}

func biasJanken() Hand {
	var handProbabilities = map[Hand]int{
		ROCK:       26,
		SCISSORS:   32,
		PAPER:      37,
		INVINCIBLE: 5,
	}

	// Calculation of the sum of all probabilities
	var allProb = 0
	for _, probability := range handProbabilities {
		allProb += probability
	}

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(allProb)

	var cumulativeProb = 0
	for hand, probability := range handProbabilities {
		cumulativeProb += probability
		if randNum < cumulativeProb {
			return hand
		}
	}
	if allProb < cumulativeProb {
		log.Fatal("Invalid prob")
	}
	panic("panic")
}

func doJanken(playerHand Hand, yodogawaHand Hand) Outcome {
	if playerHand == yodogawaHand {
		return DRAW
	} else if (playerHand == ROCK && yodogawaHand == SCISSORS) ||
		(playerHand == SCISSORS && yodogawaHand == PAPER) ||
		(playerHand == PAPER && yodogawaHand == ROCK) {
		return WIN
	} else if playerHand == HLOVE {
		return LOVE
	} else if playerHand == HHANDSHAKE {
		return HANDSHAKE
	}

	return LOSE
}
