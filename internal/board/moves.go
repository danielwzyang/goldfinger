package board

var (
	BCastleKS = true // black can castle kingside until rook or king moves
	BCastleQS = true // black can castle queenside until rook or king moves
	WCastleKS = true // white can castle kingside until rook or king moves
	WCastleQS = true // white can castle queenside until rook or king moves
)

func KingsideCastle(color byte) {
	row := 7
	if color == 'b' {
		row = 0
	}

	MakeMove(row, 4, row, 6, false) // move king to g
	MakeMove(row, 7, row, 5, false) // move rook to f
}

func QueensideCastle(color byte) {
	row := 7
	if color == 'b' {
		row = 0
	}

	MakeMove(row, 4, row, 2, false) // move king to c
	MakeMove(row, 0, row, 3, false) // move rook to d
}

// validation is only true when moves have to be made to check for something e.g. checking for checks when castling
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
