package board

func GetAllValidMoves(color byte) ([][2][2]int, int) {
	valid := [][2][2]int{}
	n := 0

	for r, row := range Board {
		for c, piece := range row {
			if piece[0] == color {
				for _, move := range getValidMoves(r, c) {
					if IsValidMove(r, c, move[0], move[1]) {
						valid = append(valid, [2][2]int{{r, c}, move})
						n++
					}
				}
			}
		}
	}

	return valid, n
}

func GetPossiblePieces(color byte, piece byte, rf int, cf int) [][2]int {
	possiblePieces := [][2]int{}
	for r, x := range Board {
		for c, y := range x {
			if y[0] == color && y[1] == piece && IsValidMove(r, c, rf, cf) {
				possiblePieces = append(possiblePieces, [2]int{r, c})
			}
		}
	}

	return possiblePieces
}

func ContainsPosition(positions [][2]int, position [2]int) bool {
	for _, item := range positions {
		if item[0] == position[0] && item[1] == position[1] {
			return true
		}
	}

	return false
}

func IsValidMove(r1 int, c1 int, r2 int, c2 int) bool {
	return ContainsPosition(getValidMoves(r1, c1), [2]int{r2, c2}) && moveOutOfCheck(r1, c1, r2, c2)
}

func ValidateMoves(r int, c int, moves [][2]int) [][2]int {
	valid := [][2]int{}
	for _, move := range moves {
		if IsValidMove(r, c, move[0], move[1]) {
			valid = append(valid, move)
		}
	}

	return valid
}

func getValidMoves(r int, c int) [][2]int {
	piece := Board[r][c]

	switch piece[1] {
	case 'P':
		return getPawnMoves(piece[0], r, c)
	case 'N':
		return getKnightMoves(piece[0], r, c)
	case 'B':
		return getBishopMoves(piece[0], r, c)
	case 'R':
		return getRookMoves(piece[0], r, c)
	case 'Q':
		return getQueenMoves(piece[0], r, c)
	case 'K':
		return getKingMoves(piece[0], r, c)
	}

	return [][2]int{}
}

func moveOutOfCheck(r1 int, c1 int, r2 int, c2 int) bool {
	color := Board[r1][c1][0]
	temp := Board[r2][c2]

	Board[r2][c2] = Board[r1][c1]
	Board[r1][c1] = " "

	if Board[r2][c2][1] == 'K' {
		if Board[r2][c2][0] == 'w' {
			WhiteKing = [2]int{r2, c2}
		} else {
			BlackKing = [2]int{r2, c2}
		}
	}

	outOfCheck := !InCheck(color)

	if Board[r2][c2][1] == 'K' {
		if Board[r2][c2][0] == 'w' {
			WhiteKing = [2]int{r1, c1}
		} else {
			BlackKing = [2]int{r1, c1}
		}
	}

	Board[r1][c1] = Board[r2][c2]
	Board[r2][c2] = temp

	return outOfCheck
}

func InCheck(color byte) bool {
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

	possiblePawnAttacks := getPawnMoves(color, king[0], king[1])
	for _, pos := range possiblePawnAttacks {
		if Board[pos[0]][pos[1]] == enemy+"P" {
			return true
		}
	}

	possibleKnightAttacks := getKnightMoves(color, king[0], king[1])
	for _, pos := range possibleKnightAttacks {
		if Board[pos[0]][pos[1]] == enemy+"N" {
			return true
		}
	}

	possibleBishopAttacks := getBishopMoves(color, king[0], king[1])
	for _, pos := range possibleBishopAttacks {
		if Board[pos[0]][pos[1]] == enemy+"B" || Board[pos[0]][pos[1]] == enemy+"Q" {
			return true
		}
	}

	possibleRookAttacks := getRookMoves(color, king[0], king[1])
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

func ValidKingSideCastle(color byte) bool {
	if InCheck(color) {
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

	// check right one
	if !moveOutOfCheck(row, 4, row, 5) {
		return false
	}

	// check right two
	if !moveOutOfCheck(row, 4, row, 6) {
		return false
	}

	return true
}

func ValidQueenSideCastle(color byte) bool {
	if InCheck(color) {
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

	// check left one
	if !moveOutOfCheck(row, 4, row, 3) {
		return false
	}

	// check left two
	if !moveOutOfCheck(row, 4, row, 2) {
		return false
	}

	return true
}
