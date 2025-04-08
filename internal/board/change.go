package board

import "math"

func KingsideCastle(color byte) {
	row := 7
	if color == 'b' {
		row = 0
	}

	MakeMove(row, 4, row, 6) // move king to g
	MakeMove(row, 7, row, 5) // move rook to f
}

func QueensideCastle(color byte) {
	row := 7
	if color == 'b' {
		row = 0
	}

	MakeMove(row, 4, row, 2) // move king to c
	MakeMove(row, 0, row, 3) // move rook to d
}

func MakeMove(r1 int, c1 int, r2 int, c2 int) {
	// update new position
	piece := Board[r1][c1]
	Board[r2][c2] = piece
	Board[r1][c1] = " "

	// if moved king, update stored king position and invalidate castling
	if piece[1] == 'K' {
		if piece[0] == 'b' {
			BlackKing = [2]int{r2, c2}
			BCastleKS = false
			BCastleQS = false
		} else {
			WhiteKing = [2]int{r2, c2}
			WCastleKS = false
			WCastleQS = false
		}
	}

	// if rook, invalidate castling for that side
	if piece[1] == 'R' {
		if piece[0] == 'b' {
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

	// if pawn handle en passant related things
	if piece[1] == 'P' {
		// pawn made double move
		if math.Abs(float64(r2-r1)) == 2 {
			// set en passant target to the space skipped over
			EnPassant = [2]int{(r1 + r2) / 2, c2}
			return
		} else if math.Abs(float64(r2-r1)) == 1 && math.Abs(float64(c2-c1)) == 1 && r2 == EnPassant[0] && c2 == EnPassant[1] {
			// capturing en passant
			// moving diagonally so distance for rows and columns are both 1 and the space is the space for en passant

			var capturedPawnRow int
			if piece[0] == 'w' {
				// white pawn moves up so the captured pawn is one row below
				capturedPawnRow = r2 + 1
			} else {
				// black pawn moving down so the captured pawn is one row above
				capturedPawnRow = r2 - 1
			}

			// capture pawn
			Board[capturedPawnRow][c2] = " "

			// reset en passant target
			EnPassant = [2]int{-10, -10}
			return
		}
	}

	// any move that isn't double pawn invalidates en passant
	EnPassant = [2]int{-10, -10}
}

func IsCapture(move [2][2]int) bool {
	return Board[move[1][0]][move[1][1]] != " " && Board[move[1][0]][move[1][1]][0] != Board[move[0][0]][move[0][1]][0]
}

func GetCaptureMoves(color byte) [][2][2]int {
	captures := [][2][2]int{}

	for r, row := range Board {
		for c, piece := range row {
			if piece[0] == color {
				for _, move := range getValidMoves(r, c) {
					if IsCapture([2][2]int{{r, c}, move}) {
						captures = append(captures, [2][2]int{{r, c}, move})
					}
				}
			}
		}
	}

	return captures
}

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
