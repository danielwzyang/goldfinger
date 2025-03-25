package main

func isValidSquare(color string, r int, c int) bool {
	// out of bounds
	if !(r >= 0 && r < 8 && c >= 0 && c < 8) {
		return false
	}

	// empty or capture
	if color == "b" {
		return Board[r][c] == " " || Board[r][c][:1] == "w"
	} else {
		return Board[r][c] == " " || Board[r][c][:1] == "b"
	}
}

func GetValidMoves(r int, c int) [][2]int {
	piece := Board[r][c]

	switch piece[1:] {
	case "P":
		return getPawnMoves(piece[:1], r, c)
	case "N":
		return getKnightMoves(piece[:1], r, c)
	case "B":
		return getBishopMoves(piece[:1], r, c)
	case "R":
		return getRookMoves(piece[:1], r, c)
	case "Q":
		return getQueenMoves(piece[:1], r, c)
	case "K":
		return getKingMoves(piece[:1], r, c)
	}

	return [][2]int{}
}

func getPawnMoves(color string, r int, c int) [][2]int {
	moves := [][2]int{}

	if color == "w" {
		// move forward
		if Board[r-1][c] == " " {
			moves = append(moves, [2]int{r - 1, c})
		}

		// capture left
		if Board[r-1][c-1][:1] == "b" {
			moves = append(moves, [2]int{r - 1, c - 1})
		}

		// capture right
		if Board[r-1][c+1][:1] == "b" {
			moves = append(moves, [2]int{r - 1, c + 1})
		}
	} else {
		// move forward
		if Board[r+1][c] == " " {
			moves = append(moves, [2]int{r + 1, c})
		}

		// capture left
		if Board[r+1][c-1][:1] == "w" {
			moves = append(moves, [2]int{r + 1, c - 1})
		}

		// capture right
		if Board[r+1][c+1][:1] == "w" {
			moves = append(moves, [2]int{r + 1, c + 1})
		}
	}

	return moves
}

func getKnightMoves(color string, r int, c int) [][2]int {
	knightMoves := [][2]int{
		{2, 1},
		{1, 2},
		{-2, 1},
		{-1, 2},
		{2, -1},
		{1, -2},
		{-2, -1},
		{-1, -2},
	}

	moves := [][2]int{}

	for _, move := range knightMoves {
		if isValidSquare(color, r+move[0], c+move[1]) {
			moves = append(moves, [2]int{r + move[0], c + move[1]})
		}
	}

	return moves
}

func getRookMoves(color string, r int, c int) [][2]int {
	moves := [][2]int{}

	// upwards
	for i := r - 1; i >= 0; i-- {
		if isValidSquare(color, i, c) {
			break
		}
		moves = append(moves, [2]int{i, c})
	}

	// downwards
	for i := r + 1; i < 8; i++ {
		if isValidSquare(color, i, c) {
			break
		}
		moves = append(moves, [2]int{i, c})
	}

	// leftwards
	for i := c - 1; i >= 0; i-- {
		if isValidSquare(color, r, i) {
			break
		}
		moves = append(moves, [2]int{r, i})
	}

	// rightwards
	for i := c + 1; i < 8; i++ {
		if isValidSquare(color, r, i) {
			break
		}
		moves = append(moves, [2]int{r, i})
	}

	return moves
}

func getBishopMoves(color string, r int, c int) [][2]int {
	moves := [][2]int{}

	// top left
	for i := 1; i < 8; i++ {
		if isValidSquare(color, r-i, c-i) {
			break
		}
		moves = append(moves, [2]int{r - i, c - i})
	}

	// top right
	for i := 1; i < 8; i++ {
		if isValidSquare(color, r-i, c+i) {
			break
		}
		moves = append(moves, [2]int{r - i, c + i})
	}

	// bottom left
	for i := 1; i < 8; i++ {
		if isValidSquare(color, r+i, c-i) {
			break
		}
		moves = append(moves, [2]int{r + i, c - i})
	}

	// bottom right
	for i := 1; i < 8; i++ {
		if isValidSquare(color, r+i, c+i) {
			break
		}
		moves = append(moves, [2]int{r + i, c + i})
	}

	return moves
}

func getQueenMoves(color string, r int, c int) [][2]int {
	return append(getRookMoves(color, r, c), getBishopMoves(color, r, c)...)
}

func getKingMoves(color string, r int, c int) [][2]int {
	kingMoves := [][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	moves := [][2]int{}

	for _, move := range kingMoves {
		if isValidSquare(color, r+move[0], c+move[1]) {
			moves = append(moves, [2]int{r + move[0], c + move[1]})
		}
	}

	return moves
}
