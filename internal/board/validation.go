package board

import (
	"regexp"
)

func ValidInput(move string, player byte) (bool, string) {
	// basic notation for moves
	regexPattern := "^([a-h][1-8])-([a-h][1-8])$"

	regex := regexp.MustCompile(regexPattern)
	if !regex.MatchString(move) {
		return false, "Please follow the format described above. Read /docs/inputs.md for help."
	}

	// turning the moves into matrix coordinates
	alpha1 := move[:2]
	alpha2 := move[3:]

	p1 := alphaToNumeric(alpha1)
	p2 := alphaToNumeric(alpha2)

	// if the piece isn't the player's
	if Board[p1[0]][p1[1]][0] != player {
		return false, "You do not own the piece at " + alpha1 + "."
	}

	// the position it's trying to move isn't valid
	if !isValidMove(p1[0], p1[1], p2[0], p2[1]) {
		return false, "The piece at " + alpha1 + " cannot move to " + alpha2 + "."
	}

	return true, ""
}

func isValidMove(r1 int, c1 int, r2 int, c2 int) bool {
	validMoves := getValidMoves(r1, c1)

	for _, move := range validMoves {
		if move[0] == r2 && move[1] == c2 {
			return true
		}
	}

	return false
}

func getValidMoves(r int, c int) [][2]int {
	piece := Board[r][c]

	switch piece[1:] {
	case "P":
		return getPawnMoves(piece[0], r, c)
	case "N":
		return getKnightMoves(piece[0], r, c)
	case "B":
		return getBishopMoves(piece[0], r, c)
	case "R":
		return getRookMoves(piece[0], r, c)
	case "Q":
		return getQueenMoves(piece[0], r, c)
	case "K":
		return getKingMoves(piece[0], r, c)
	}

	return [][2]int{}
}

func outOfBounds(r int, c int) bool {
	return !(r >= 0 && r < 8 && c >= 0 && c < 8)
}

func isValidSquare(color byte, r int, c int) bool {
	if outOfBounds(r, c) {
		return false
	}

	// empty or capture
	if color == 'b' {
		return Board[r][c] == " " || Board[r][c][0] == 'w'
	} else {
		return Board[r][c] == " " || Board[r][c][0] == 'b'
	}
}

func getPawnMoves(color byte, r int, c int) [][2]int {
	moves := [][2]int{}

	if color == 'w' {
		// move forward
		if !outOfBounds(r-1, c) && Board[r-1][c] == " " {
			moves = append(moves, [2]int{r - 1, c})
		}

		// capture left
		if !outOfBounds(r-1, c-1) && Board[r-1][c-1][0] == 'b' {
			moves = append(moves, [2]int{r - 1, c - 1})
		}

		// capture right
		if !outOfBounds(r-1, c+1) && Board[r-1][c+1][0] == 'b' {
			moves = append(moves, [2]int{r - 1, c + 1})
		}
	} else {
		// move forward
		if !outOfBounds(r+1, c) && Board[r+1][c] == " " {
			moves = append(moves, [2]int{r + 1, c})
		}

		// capture left
		if !outOfBounds(r+1, c-1) && Board[r+1][c-1][0] == 'w' {
			moves = append(moves, [2]int{r + 1, c - 1})
		}

		// capture right
		if !outOfBounds(r+1, c+1) && Board[r+1][c+1][0] == 'w' {
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
	}

	// downwards
	for i := r + 1; i < 8; i++ {
		if !isValidSquare(color, i, c) {
			break
		}
		moves = append(moves, [2]int{i, c})
	}

	// leftwards
	for i := c - 1; i >= 0; i-- {
		if !isValidSquare(color, r, i) {
			break
		}
		moves = append(moves, [2]int{r, i})
	}

	// rightwards
	for i := c + 1; i < 8; i++ {
		if !isValidSquare(color, r, i) {
			break
		}
		moves = append(moves, [2]int{r, i})
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
	}

	// top right
	for i := 1; i < 8; i++ {
		if !isValidSquare(color, r-i, c+i) {
			break
		}
		moves = append(moves, [2]int{r - i, c + i})
	}

	// bottom left
	for i := 1; i < 8; i++ {
		if !isValidSquare(color, r+i, c-i) {
			break
		}
		moves = append(moves, [2]int{r + i, c - i})
	}

	// bottom right
	for i := 1; i < 8; i++ {
		if !isValidSquare(color, r+i, c+i) {
			break
		}
		moves = append(moves, [2]int{r + i, c + i})
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
		if !isValidSquare(color, r+move[0], c+move[1]) {
			moves = append(moves, [2]int{r + move[0], c + move[1]})
		}
	}

	return moves
}
