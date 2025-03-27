package board

var (
	BCastleKS = true // black can castle kingside until rook or king moves
	BCastleQS = true // black can castle queenside until rook or king moves
	WCastleKS = true // white can castle kingside until rook or king moves
	WCastleQS = true // white can castle queenside until rook or king moves
)

func alphaToNumeric(position string) [2]int {
	// converts string to rune array
	runes := []rune(position)

	// ascii value of 'a' is 97, '8' is 56
	return [2]int{56 - int(runes[1]), int(runes[0]) - 97}
}

// assumes moves are already valid
func ParseMove(move string, color byte) {
	// handle castling notation
	if move == "0-0" {
		row := 7
		if color == 'b' {
			row = 0
		}

		MakeMove(row, 4, row, 6, false) // move king to g
		MakeMove(row, 7, row, 5, false) // move rook to f

		return
	}

	if move == "0-0-0" {
		row := 7
		if color == 'b' {
			row = 0
		}

		MakeMove(row, 4, row, 2, false) // move king to c
		MakeMove(row, 0, row, 3, false) // move rook to d

		return
	}

	// convert to matrix coordinates
	alpha1 := move[:2]
	alpha2 := move[3:]

	p1 := alphaToNumeric(alpha1)
	p2 := alphaToNumeric(alpha2)

	// make move
	MakeMove(p1[0], p1[1], p2[0], p2[1], false)
}

func MakeMove(r1 int, c1 int, r2 int, c2 int, validation bool) {
	// update new position
	Board[r2][c2] = Board[r1][c1]

	// update original position
	Board[r1][c1] = " "

	// if moved king, update stored king position and invalidate castling
	if Board[r2][c2][1] == 'K' {
		if Board[r2][c2][0] == 'b' {
			BlackKing = [2]int{r2, c2}

			if !validation {
				BCastleKS = false
				BCastleQS = false
			}
		} else {
			WhiteKing = [2]int{r2, c2}

			if !validation {
				WCastleKS = false
				WCastleQS = false
			}
		}
	}

	// if rook, invalidate castling
	if Board[r2][c2][1] == 'R' && !validation {
		if Board[r2][c2][0] == 'b' {
			if c1 == 7 {
				BCastleKS = false
			}
			if c1 == 0 {
				BCastleQS = false
			}
		} else {
			if c1 == 7 {
				WCastleKS = false
			}
			if c1 == 0 {
				WCastleQS = false
			}
		}
	}
}
