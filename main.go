package main

import "fmt"

func main() {
	var inputHand string
	fmt.Print("Input: ")
	fmt.Scanf("%s", &inputHand)

	playerHand := getPlayerHand(inputHand)
	yodogawaHand := biasJanken()

	fmt.Println("Your hand: ", playerHand)
	fmt.Println("Yodogawa-san hand: ", yodogawaHand)

	result := doJanken(playerHand, yodogawaHand)
	fmt.Println(result)
}
