package board

func getPawnMoves(from Position) []Move {
	r := from.Rank
	c := from.File

	moves := []Move{}

	if Board[r][c].Color == WHITE {
		// move forward one
		if r != 0 && Board[r-1][c].Type == EMPTY {
			moves = append(moves, Move{from, Position{r - 1, c}})
		}

		// move forward two from the starting row (row 6 for white)
		if r == 6 && Board[r-1][c].Type == EMPTY && Board[r-2][c].Type == EMPTY {
			moves = append(moves, Move{from, Position{r - 2, c}})
		}

		// capture left
		if r != 0 && c != 0 {
			// regular capture
			if Board[r-1][c-1].Color == BLACK {
				moves = append(moves, Move{from, Position{r - 1, c - 1}})
				// en passant capture
			} else if EnPassant.Rank == r-1 && EnPassant.File == c-1 {
				moves = append(moves, Move{from, Position{r - 1, c - 1}})
			}
		}

		// capture right
		if r != 0 && c != 7 {
			// regular capture
			if Board[r-1][c+1].Color == BLACK {
				moves = append(moves, Move{from, Position{r - 1, c + 1}})
				// en passant capture
			} else if EnPassant.Rank == r-1 && EnPassant.File == c+1 {
				moves = append(moves, Move{from, Position{r - 1, c + 1}})
			}
		}
	} else if Board[r][c].Color == BLACK {
		// move forward one
		if r != 7 && Board[r+1][c].Type == EMPTY {
			moves = append(moves, Move{from, Position{r + 1, c}})
		}

		// move forward two from the starting row (row 1 for black)
		if r == 1 && Board[r+1][c].Type == EMPTY && Board[r+2][c].Type == EMPTY {
			moves = append(moves, Move{from, Position{r + 2, c}})
		}

		// capture left
		if r != 7 && c != 0 {
			// regular capture
			if Board[r+1][c-1].Color == WHITE {
				moves = append(moves, Move{from, Position{r + 1, c - 1}})
				// en passant capture
			} else if EnPassant.Rank == r+1 && EnPassant.File == c-1 {
				moves = append(moves, Move{from, Position{r + 1, c - 1}})
			}
		}

		// capture right
		if r != 7 && c != 7 {
			// regular capture
			if Board[r+1][c+1].Color == WHITE {
				moves = append(moves, Move{from, Position{r + 1, c + 1}})
				// en passant capture
			} else if EnPassant.Rank == r+1 && EnPassant.File == c+1 {
				moves = append(moves, Move{from, Position{r + 1, c + 1}})
			}
		}
	}

	return moves
}

func getKnightMoves(from Position) []Move {
	r := from.Rank
	c := from.File

	attacking := []Position{
		{2, 1},
		{1, 2},
		{-2, 1},
		{-1, 2},
		{2, -1},
		{1, -2},
		{-2, -1},
		{-1, -2},
	}

	moves := []Move{}

	for _, position := range attacking {
		to := Position{r + position.Rank, c + position.File}

		if isValidPosition(to, Board[r][c].Color) {
			moves = append(moves, Move{from, to})
		}
	}

	return moves
}

func getBishopMoves(from Position) []Move {
	r := from.Rank
	c := from.File

	color := Board[r][c].Color

	moves := []Move{}

	// top left
	for i := 1; i < 8; i++ {
		to := Position{r - i, c - i}
		if !isValidPosition(to, color) {
			break
		}

		moves = append(moves, Move{from, to})

		// capture
		if Board[r-i][c-i].Color == color^1 {
			break
		}
	}

	// top right
	for i := 1; i < 8; i++ {
		to := Position{r - i, c + i}

		if !isValidPosition(to, color) {
			break
		}

		moves = append(moves, Move{from, to})

		// capture
		if Board[r-i][c+i].Color == color^1 {
			break
		}
	}

	// bottom left
	for i := 1; i < 8; i++ {
		to := Position{r + i, c - i}

		if !isValidPosition(to, color) {
			break
		}

		moves = append(moves, Move{from, to})

		// capture
		if Board[r+i][c-i].Color == color^1 {
			break
		}
	}

	// bottom right
	for i := 1; i < 8; i++ {
		to := Position{r + i, c + i}

		if !isValidPosition(to, color) {
			break
		}

		moves = append(moves, Move{from, to})

		// capture
		if Board[r+i][c+i].Color == color^1 {
			break
		}
	}

	return moves
}

func getRookMoves(from Position) []Move {
	r := from.Rank
	c := from.File

	color := Board[r][c].Color

	moves := []Move{}

	// upwards
	for i := r - 1; i >= 0; i-- {
		to := Position{i, c}

		if !isValidPosition(to, color) {
			break
		}

		moves = append(moves, Move{from, to})

		// capture
		if Board[i][c].Color == color^1 {
			break
		}
	}

	// downwards
	for i := r + 1; i < 8; i++ {
		to := Position{i, c}

		if !isValidPosition(to, color) {
			break
		}

		moves = append(moves, Move{from, to})

		// capture
		if Board[i][c].Color == color^1 {
			break
		}
	}

	// leftwards
	for i := c - 1; i >= 0; i-- {
		to := Position{r, i}

		if !isValidPosition(to, color) {
			break
		}

		moves = append(moves, Move{from, to})

		// capture
		if Board[r][i].Color == color^1 {
			break
		}
	}

	// rightwards
	for i := c + 1; i < 8; i++ {
		to := Position{r, i}

		if !isValidPosition(to, color) {
			break
		}

		moves = append(moves, Move{from, to})

		// capture
		if Board[r][i].Color == color^1 {
			break
		}
	}

	return moves
}

func getQueenMoves(from Position) []Move {
	return append(getRookMoves(from), getBishopMoves(from)...)
}

func getKingMoves(from Position) []Move {
	r := from.Rank
	c := from.File

	attacking := []Position{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	moves := []Move{}

	for _, position := range attacking {
		to := Position{r + position.Rank, c + position.File}

		if isValidPosition(to, Board[r][c].Color) {
			moves = append(moves, Move{from, to})
		}
	}

	return moves
}
