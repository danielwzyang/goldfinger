package main

import "fmt"

func main() {
	InitBoard()

	ClearTerminal()

	PrintBoard()

	fmt.Println(GetValidMoves(0, 1))
}

func ClearTerminal() {
	fmt.Print("\033[H\033[2J")
}
