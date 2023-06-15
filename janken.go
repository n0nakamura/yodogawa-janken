package main

import (
	"errors"
	"log"
	"math/rand"
	"regexp"
	"time"
)

type Outcome uint

const (
	WIN Outcome = iota
	LOSE
	DRAW
	LOVE
	HANDSHAKE
	INFO
	BATTLE
)

var outcomeNameMap = map[Outcome]string{
	WIN:       "YOU WIN",
	LOSE:      "YOU LOSE",
	DRAW:      "DRAW",
	LOVE:      "BIG LOVE...🤟",
	HANDSHAKE: "SHAKING... 🤝",
	INFO:      name + "\nversion " + version,
	BATTLE:    "Not ready yet...",
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

var handPatterns = map[Hand]string{
	ROCK:       `R✊👊🤛🤜💪🪨`,
	SCISSORS:   `S✌🤞🦞🦀🦂✂︎✃✄💇💇‍♂️💇‍♀️`,
	PAPER:      `P🖐✋🤚🖖🫲🫱🫳🫴👋👐🤲🤗🪬🧻📝📄📃📜📑🧾📰🗺️🧧🔖🗞️🙋🙋‍♂️🙋‍♀️`,
	HLOVE:      `🤟🫶🫂`,
	HHANDSHAKE: `🤝`,
	OTHER:      `👌🤌🤏🤘🤙👈👉👆👇☝👍👏🙏🫵`,
}

func (h Hand) String() string {
	if s, ok := handNames[h]; ok {
		return s
	}
	return "Unknown hand"
}

type Mode uint

const (
	M_INFO Mode = iota
	M_BATTLE
)

var modeNames = map[Mode]string{
	M_INFO:   "ℹ️ Info",
	M_BATTLE: "⚔️ Battle",
}

var modePatterns = map[Mode]string{
	M_INFO:   `Iℹ️`,
	M_BATTLE: "B⚔️",
}

func getPlayerHand(playerHand string) (Hand, error) {
	var re = make(map[Hand]*regexp.Regexp)
	for hand, pattern := range handPatterns {
		re[hand] = regexp.MustCompile(`[` + pattern + `]`)
		if ok := re[hand].MatchString(playerHand); ok {
			return hand, nil
		}
	}
	return 0, errors.New("invalid hand")
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
	switch {
	case playerHand == yodogawaHand:
		return DRAW
	case (playerHand == ROCK && yodogawaHand == SCISSORS) ||
		(playerHand == SCISSORS && yodogawaHand == PAPER) ||
		(playerHand == PAPER && yodogawaHand == ROCK):
		return WIN
	case playerHand == HLOVE:
		return LOVE
	case playerHand == HHANDSHAKE:
		return HANDSHAKE
	default:
		return LOSE
	}
}
