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
	M_LTW: {
		InputPattern: P_LTW,
		DoFunc:       ltw,
	},
	M_OMIKUJI: {
		InputPattern: P_OMIKUJI,
		DoFunc:       omikuji,
	},
	M_EMOJI: {
		InputPattern: P_EMOJI,
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
	M_LTW
	M_OMIKUJI
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
	P_PAPER     = `P🖐✋🤚🖖🫲🫱🫳🫴🫸🫷👋👐🤲🤗🪬🧻📝📄📃📜📑🧾📰🗺️🧧🔖🗞️🙋🙋‍♂️🙋‍♀️`
	P_JANKEN    = P_ROCK + P_SCISSORS + P_PAPER
	P_UP        = `U👆☝`
	P_DOWN      = `D👇`
	P_RIGHT     = `R👉`
	P_LEFT      = `L👈`
	P_FRONT     = `F🫵`
	P_LTW       = P_UP + P_DOWN + P_RIGHT + P_LEFT + P_FRONT
	P_OMIKUJI   = `O👏🙏`
	P_LOVE      = `🤟🫶🫂`
	P_SHAKE     = `🤝`
	P_OTHERHAND = `👌🤌🤏🤘🤙👍`
	P_BAMBOO    = `🎍`
	P_BROCCOLI  = `🥦`
	P_COCHLEA   = `🐌`
	P_EMOJI     = P_LOVE +
		P_SHAKE +
		P_OTHERHAND +
		P_BAMBOO +
		P_BROCCOLI +
		P_COCHLEA
	P_INFO   = `Iℹ️`
	P_BATTLE = `B⚔️`
)

var ErrNoValuesIncluded = errors.New("contains no matching values")

func generateContent(pcontent string) (string, error) {
	var re = make(map[ModeID]*regexp.Regexp)
	for mid, m := range modes {
		re[mid] = regexp.MustCompile(`[` + m.InputPattern + `]`)
		if ok := re[mid].MatchString(pcontent); ok {
			return m.DoFunc(pcontent)
		}
	}
	return "", ErrNoValuesIncluded
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
			break
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

func ltw(pcontent string) (string, error) {
	type Hand uint
	const (
		UP Hand = iota
		DOWN
		RIGHT
		LEFT
		FRONT
		NUM // Total number
	)
	var handNames = map[Hand]string{
		UP:    "👆 Up",
		DOWN:  "👇 Down",
		RIGHT: "👉 Right",
		LEFT:  "👈 Left",
		FRONT: "🫵 Front",
	}
	var handPatterns = map[Hand]string{
		UP:    P_UP,
		DOWN:  P_DOWN,
		RIGHT: P_RIGHT,
		LEFT:  P_LEFT,
		FRONT: P_FRONT,
	}

	type Result uint
	const (
		WIN Result = iota
		LOSE
	)
	var resultNameMap = map[Result]string{
		WIN:  "YOU WIN",
		LOSE: "YOU LOSE",
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

	rand.Seed(time.Now().UnixNano())
	yodogawaHand := Hand(rand.Intn(int(NUM)))

	// Get result
	var result Result
	switch {
	case playerHand == yodogawaHand:
		result = WIN
	default:
		result = LOSE
	}

	return "Your hand: " + handNames[playerHand] + "\n" +
		"Yodogawa-san hand: " + handNames[yodogawaHand] + "\n" +
		resultNameMap[result], nil
}

func omikuji(pcontent string) (string, error) {
	f := []string{
		"大吉",
		"吉",
		"中吉",
		"小吉",
		"半吉",
		"末吉",
		"末小吉",
		"平",
		"凶",
		"小凶",
		"半凶",
		"末凶",
		"大凶",
	}

	rand.Seed(time.Now().UnixNano())
	p := rand.Intn(len(f))

	return "⛩️ Your fortune is " + f[p], nil
}

func emoji(pcontent string) (string, error) {
	const (
		LOVE = iota
		SHAKE
		OTHERHAND
		BAMBOO
		BROCCOLI
		COCHLEA
	)
	var pattern = []string{
		P_LOVE,
		P_SHAKE,
		P_OTHERHAND,
		P_BAMBOO,
		P_BROCCOLI,
		P_COCHLEA,
	}

	// Get player hand
	var re = make(map[int]*regexp.Regexp)
	for i, p := range pattern {
		re[i] = regexp.MustCompile(`[` + p + `]`)
	}

	switch {
	case re[LOVE].MatchString(pcontent):
		return "🤟 BIG LOVE...", nil
	case re[SHAKE].MatchString(pcontent):
		return "🤝 SHAKING...", nil
	case re[OTHERHAND].MatchString(pcontent):
		return re[OTHERHAND].FindString(pcontent), nil
	case re[BAMBOO].MatchString(pcontent):
		return "🎍 これは竹", nil
	case re[BROCCOLI].MatchString(pcontent):
		return "🥦 https://cookpad.com/search/%E3%83%96%E3%83%AD%E3%83%83%E3%82%B3%E3%83%AA%E3%83%BC", nil
	case re[COCHLEA].MatchString(pcontent):
		return "\n₍₍🐌⁾⁾\n\n見て！カタツムリが踊っているよ\nかわいいね\n\n₍₍⁽⁽🐌₎₎⁾⁾\n\nみんながYodogawa-Jankenに反応してくれるので、カタツムリはさらに踊りだしました\nあなたのおかげです\nありがとう", nil
	default:
		return "🤔", nil
	}
}

func info(pcontent string) (string, error) {
	return name + " " + version, nil
}

func battle(pcontent string) (string, error) {
	return "Not available yet...", nil
}
