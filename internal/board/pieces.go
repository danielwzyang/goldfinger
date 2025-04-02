package board

func getPawnMoves(color byte, r int, c int) [][2]int {
	moves := [][2]int{}

	if color == 'w' {
		// move forward
		if Board[r-1][c] == " " {
			moves = append(moves, [2]int{r - 1, c})
		}

		// move forward two
		if r == 6 && Board[r-2][c] == " " {
			moves = append(moves, [2]int{r - 2, c})
		}

		// capture left
		if c != 0 && Board[r-1][c-1][0] == 'b' {
			moves = append(moves, [2]int{r - 1, c - 1})
		}

		// capture right
		if c != 7 && Board[r-1][c+1][0] == 'b' {
			moves = append(moves, [2]int{r - 1, c + 1})
		}
	} else {
		// move forward
		if Board[r+1][c] == " " {
			moves = append(moves, [2]int{r + 1, c})
		}

		// move forward two
		if r == 1 && Board[r+2][c] == " " {
			moves = append(moves, [2]int{r + 2, c})
		}

		// capture left
		if c != 0 && Board[r+1][c-1][0] == 'w' {
			moves = append(moves, [2]int{r + 1, c - 1})
		}

		// capture right
		if c != 7 && Board[r+1][c+1][0] == 'w' {
			moves = append(moves, [2]int{r + 1, c + 1})
		}
	}

	return moves
}

func getKnightMoves(color byte, r int, c int) [][2]int {
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

func getRookMoves(color byte, r int, c int) [][2]int {
	moves := [][2]int{}

	// upwards
	for i := r - 1; i >= 0; i-- {
		if !isValidSquare(color, i, c) {
			break
		}

		moves = append(moves, [2]int{i, c})

		// capture
		if Board[i][c][0] != ' ' && Board[i][c][0] != color {
			break
		}
	}

	// downwards
	for i := r + 1; i < 8; i++ {
		if !isValidSquare(color, i, c) {
			break
		}

		moves = append(moves, [2]int{i, c})

		// capture
		if Board[i][c][0] != ' ' && Board[i][c][0] != color {
			break
		}
	}

	// leftwards
	for i := c - 1; i >= 0; i-- {
		if !isValidSquare(color, r, i) {
			break
		}

		moves = append(moves, [2]int{r, i})

		// capture
		if Board[r][i][0] != ' ' && Board[r][i][0] != color {
			break
		}
	}

	// rightwards
	for i := c + 1; i < 8; i++ {
		if !isValidSquare(color, r, i) {
			break
		}

		moves = append(moves, [2]int{r, i})

		// capture
		if Board[r][i][0] != ' ' && Board[r][i][0] != color {
			break
		}
	}

	return moves
}

func getBishopMoves(color byte, r int, c int) [][2]int {
	moves := [][2]int{}

	// top left
	for i := 1; i < 8; i++ {
		if !isValidSquare(color, r-i, c-i) {
			break
		}

		moves = append(moves, [2]int{r - i, c - i})

		// capture
		if Board[r-i][c-i][0] != ' ' && Board[r-i][c-i][0] != color {
			break
		}
	}

	// top right
	for i := 1; i < 8; i++ {
		if !isValidSquare(color, r-i, c+i) {
			break
		}

		moves = append(moves, [2]int{r - i, c + i})

		// capture
		if Board[r-i][c+i][0] != ' ' && Board[r-i][c+i][0] != color {
			break
		}
	}

	// bottom left
	for i := 1; i < 8; i++ {
		if !isValidSquare(color, r+i, c-i) {
			break
		}

		moves = append(moves, [2]int{r + i, c - i})

		// capture
		if Board[r+i][c-i][0] != ' ' && Board[r+i][c-i][0] != color {
			break
		}
	}

	// bottom right
	for i := 1; i < 8; i++ {
		if !isValidSquare(color, r+i, c+i) {
			break
		}

		moves = append(moves, [2]int{r + i, c + i})

		// capture
		if Board[r+i][c+i][0] != ' ' && Board[r+i][c+i][0] != color {
			break
		}
	}

	return moves
}

func getQueenMoves(color byte, r int, c int) [][2]int {
	return append(getRookMoves(color, r, c), getBishopMoves(color, r, c)...)
}

func getKingMoves(color byte, r int, c int) [][2]int {
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
