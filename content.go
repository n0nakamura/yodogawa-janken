package main

import (
	"errors"
	"log"
	"math/rand"
	"regexp"
	"time"
)

var modes = map[ModeID]S_Mode{
	M_JANKEN: {
		InputPattern: P_JANKEN,
		DoFunc:       janken,
	},
	M_EMOJI: {
		InputPattern: P_OTHERHAND,
		DoFunc:       emoji,
	},
	M_INFO: {
		InputPattern: P_INFO,
		DoFunc:       info,
	},
	M_BATTLE: {
		InputPattern: P_BATTLE,
		DoFunc:       battle,
	},
}

type ModeID uint

const (
	M_JANKEN ModeID = iota
	M_EMOJI
	M_INFO
	M_BATTLE
)

type S_Mode struct {
	InputPattern string
	DoFunc       func(pcontent string) (string, error)
}

const (
	P_ROCK      = `R✊👊🤛🤜💪🪨`
	P_SCISSORS  = `S✌🤞🦞🦀🦂✂︎✃✄💇💇‍♂️💇‍♀️`
	P_PAPER     = `P🖐✋🤚🖖🫲🫱🫳🫴👋👐🤲🤗🪬🧻📝📄📃📜📑🧾📰🗺️🧧🔖🗞️🙋🙋‍♂️🙋‍♀️`
	P_JANKEN    = P_ROCK + P_SCISSORS + P_PAPER
	P_LOVE      = `🤟🫶🫂`
	P_SHAKE     = `🤝`
	P_OTHERHAND = `👌🤌🤏🤘🤙👈👉👆👇☝👍👏🙏🫵`
	P_EMOJI     = P_LOVE + P_SHAKE + P_OTHERHAND
	P_INFO      = `Iℹ️`
	P_BATTLE    = `B⚔️`
)

var ErrNoValuesIncluded = errors.New("contains no matching values")

func generateContent(pcontent string) (string, error) {
	var re = make(map[ModeID]*regexp.Regexp)
	var mode S_Mode
	var success = 0
	for mid, m := range modes {
		re[mid] = regexp.MustCompile(`[` + m.InputPattern + `]`)
		if ok := re[mid].MatchString(pcontent); ok {
			mode = m
			success++
		}
	}
	if success == 0 {
		return "", ErrNoValuesIncluded
	}

	return mode.DoFunc(pcontent)
}

func janken(pcontent string) (string, error) {
	type Hand uint
	const (
		ROCK Hand = iota
		SCISSORS
		PAPER
		INVINCIBLE
	)
	var handNames = map[Hand]string{
		ROCK:       "✊ Rock",
		SCISSORS:   "✌ Scissors",
		PAPER:      "🖐 Paper",
		INVINCIBLE: "👉 Invincible",
	}
	var handPatterns = map[Hand]string{
		ROCK:     P_ROCK,
		SCISSORS: P_SCISSORS,
		PAPER:    P_PAPER,
	}
	var handProbabilities = map[Hand]int{
		ROCK:       26,
		SCISSORS:   32,
		PAPER:      37,
		INVINCIBLE: 5,
	}

	type Result uint
	const (
		WIN Result = iota
		LOSE
		DRAW
	)
	var resultNameMap = map[Result]string{
		WIN:  "YOU WIN",
		LOSE: "YOU LOSE",
		DRAW: "DRAW",
	}

	// Get player hand
	var re = make(map[Hand]*regexp.Regexp)
	var playerHand Hand
	for h, pattern := range handPatterns {
		re[h] = regexp.MustCompile(`[` + pattern + `]`)
		if ok := re[h].MatchString(pcontent); ok {
			playerHand = h
		}
	}
	// ここで不適切なplayerHandでエラーを返す。

	// Calculation of the sum of all probabilities
	var allProb = 0
	for _, probability := range handProbabilities {
		allProb += probability
	}

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(allProb)

	var cumulativeProb = 0
	var yodogawaHand Hand
	for h, probability := range handProbabilities {
		cumulativeProb += probability
		if randNum < cumulativeProb {
			yodogawaHand = h
		}
	}
	if allProb < cumulativeProb {
		log.Fatal("Invalid prob")
		// いい感じのエラーを返す。
	}

	// Get result
	var result Result
	switch {
	case playerHand == yodogawaHand:
		result = DRAW
	case (playerHand == ROCK && yodogawaHand == SCISSORS) ||
		(playerHand == SCISSORS && yodogawaHand == PAPER) ||
		(playerHand == PAPER && yodogawaHand == ROCK):
		result = WIN
	default:
		result = LOSE
	}

	return "Your hand: " + handNames[playerHand] + "\n" +
		"Yodogawa-san hand: " + handNames[yodogawaHand] + "\n" +
		resultNameMap[result], nil
}

func emoji(pcontent string) (string, error) {
	return "🤔", nil
}

func info(pcontent string) (string, error) {
	return name + " " + version, nil
}

func battle(pcontent string) (string, error) {
	return "Not available yet...", nil
}
