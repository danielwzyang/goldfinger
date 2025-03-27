package board

func alphaToNumeric(position string) [2]int {
	// converts string to rune array
	runes := []rune(position)

	// ascii value of 'a' is 97, '8' is 56
	return [2]int{56 - int(runes[1]), int(runes[0]) - 97}
}

func ParseMove(move string) {
	// convert to matrix coordinates
	alpha1 := move[:2]
	alpha2 := move[3:]

	p1 := alphaToNumeric(alpha1)
	p2 := alphaToNumeric(alpha2)

	// store old piece
	temp := Board[p1[0]][p1[1]]

	// update original position
	Board[p1[0]][p1[1]] = " "

	// update new position
	Board[p2[0]][p2[1]] = temp
}
