package board

import (
	"regexp"
)

func ValidInput(move string, player byte) (bool, string) {
	// castling
	if move == "0-0" {
		if !validKingSideCastle(player) {
			return false, "You cannot castle kingside."
		}
		return true, ""
	}
	if move == "0-0-0" {
		if !validQueenSideCastle(player) {
			return false, "You cannot castle queenside."
		}
		return true, ""
	}

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

	if Board[p1[0]][p1[1]][0] == ' ' {
		return false, "There is no piece at " + alpha1 + "."
	}

	// if the piece isn't the player's
	if Board[p1[0]][p1[1]][0] != player {
		return false, "You do not own the piece at " + alpha1 + "."
	}

	// the position it's trying to move isn't valid
	if !isValidMove(p1[0], p1[1], p2[0], p2[1]) {
		return false, "The piece at " + alpha1 + " cannot move to " + alpha2 + "."
	}

	MakeMove(p1[0], p1[1], p2[0], p2[1], true)

	if inCheck(player) {
		MakeMove(p2[0], p2[1], p1[0], p1[1], true)
		return false, "Your king is still in check!"
	}

	MakeMove(p2[0], p2[1], p1[0], p1[1], true)

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

	switch piece[1] {
	case 'P':
		return getValidPawnMoves(piece[0], r, c)
	case 'N':
		return getValidKnightMoves(piece[0], r, c)
	case 'B':
		return getValidBishopMoves(piece[0], r, c)
	case 'R':
		return getValidRookMoves(piece[0], r, c)
	case 'Q':
		return getValidQueenMoves(piece[0], r, c)
	case 'K':
		return getValidKingMoves(piece[0], r, c)
	}

	return [][2]int{}
}

func inCheck(color byte) bool {
	enemy := "b"
	if color == 'b' {
		enemy = "w"
	}

	king := WhiteKing
	if color == 'b' {
		king = BlackKing
	}

	// instead of going through every opponent piece to see if it can attack the king,
	// we can check enemy pieces from the king's position to save time

	possiblePawnAttacks := getValidPawnMoves(color, king[0], king[1])
	for _, pos := range possiblePawnAttacks {
		if Board[pos[0]][pos[1]] == enemy+"P" {
			return true
		}
	}

	possibleKnightAttacks := getValidKnightMoves(color, king[0], king[1])
	for _, pos := range possibleKnightAttacks {
		if Board[pos[0]][pos[1]] == enemy+"N" {
			return true
		}
	}

	possibleBishopAttacks := getValidBishopMoves(color, king[0], king[1])
	for _, pos := range possibleBishopAttacks {
		if Board[pos[0]][pos[1]] == enemy+"B" || Board[pos[0]][pos[1]] == enemy+"Q" {
			return true
		}
	}

	possibleRookAttacks := getValidRookMoves(color, king[0], king[1])
	for _, pos := range possibleRookAttacks {
		if Board[pos[0]][pos[1]] == enemy+"R" || Board[pos[0]][pos[1]] == enemy+"Q" {
			return true
		}
	}

	return false
}

func isValidSquare(color byte, r int, c int) bool {
	// out of bounds
	if !(r >= 0 && r < 8 && c >= 0 && c < 8) {
		return false
	}

	// empty or capture
	if color == 'b' {
		return Board[r][c] == " " || Board[r][c][0] == 'w'
	} else {
		return Board[r][c] == " " || Board[r][c][0] == 'b'
	}
}

func getValidPawnMoves(color byte, r int, c int) [][2]int {
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

func getValidKnightMoves(color byte, r int, c int) [][2]int {
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

func getValidRookMoves(color byte, r int, c int) [][2]int {
	moves := [][2]int{}

	// upwards
	for i := r - 1; i >= 0; i-- {
		if !isValidSquare(color, i, c) {
			break
		}

		moves = append(moves, [2]int{i, c})

		// capture
		if Board[i][c][0] != color {
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
		if Board[i][c][0] != color {
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
		if Board[r][i][0] != color {
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
		if Board[r][i][0] != color {
			break
		}
	}

	return moves
}

func getValidBishopMoves(color byte, r int, c int) [][2]int {
	moves := [][2]int{}

	// top left
	for i := 1; i < 8; i++ {
		if !isValidSquare(color, r-i, c-i) {
			break
		}

		moves = append(moves, [2]int{r - i, c - i})

		// capture
		if Board[r-i][c-i][0] != color {
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
		if Board[r-i][c+i][0] != color {
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
		if Board[r+i][c-i][0] != color {
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
		if Board[r+i][c+i][0] != color {
			break
		}
	}

	return moves
}

func getValidQueenMoves(color byte, r int, c int) [][2]int {
	return append(getValidRookMoves(color, r, c), getValidBishopMoves(color, r, c)...)
}

func getValidKingMoves(color byte, r int, c int) [][2]int {
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

func validKingSideCastle(color byte) bool {
	if inCheck(color) {
		return false
	}

	canCastle := WCastleKS
	row := 7
	if color == 'b' {
		row = 0
		canCastle = BCastleKS
	}

	// rook or king has moved or spaces aren't empty
	if !canCastle || Board[row][5] != " " || Board[row][6] != " " {
		return false
	}

	// checking if squares are in check

	// move king right one
	MakeMove(row, 4, row, 5, true)
	if inCheck(color) {
		MakeMove(row, 5, row, 4, true)
		return false
	}

	// move king right one
	MakeMove(row, 5, row, 6, true)
	if inCheck(color) {
		MakeMove(row, 6, row, 4, true)
		return false
	}

	// move king back
	MakeMove(row, 6, row, 4, true)

	return true
}

func validQueenSideCastle(color byte) bool {
	if inCheck(color) {
		return false
	}

	canCastle := WCastleQS
	row := 7
	if color == 'b' {
		row = 0
		canCastle = BCastleQS
	}

	// rook or king has moved or spaces aren't empty
	if !canCastle || Board[row][1] != " " || Board[row][2] != " " || Board[row][3] != " " {
		return false
	}

	// checking if squares are in check

	// move king left one
	MakeMove(row, 4, row, 3, true)
	if inCheck(color) {
		MakeMove(row, 3, row, 4, true)
		return false
	}

	// move king left one
	MakeMove(row, 3, row, 2, true)
	if inCheck(color) {
		MakeMove(row, 2, row, 4, true)
		return false
	}

	// move king back
	MakeMove(row, 2, row, 4, true)

	return true
}
