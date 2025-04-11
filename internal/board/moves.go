package board

func GetCaptureMoves(color int) []Move {
	captures := []Move{}

	for r, row := range Board {
		for c, piece := range row {
			if piece.Color == color {
				for _, move := range getValidMoves(Position{r, c}) {
					if IsCapture(move) {
						captures = append(captures, move)
					}
				}
			}
		}
	}

	return captures
}

func GetAllValidMoves(color int) ([]Move, int) {
	valid := []Move{}
	n := 0

	for r, row := range Board {
		for c, piece := range row {
			if piece.Color == color {
				for _, move := range getValidMoves(Position{r, c}) {
					if IsValidMove(move) {
						valid = append(valid, move)
						n++
					}
				}
			}
		}
	}

	return valid, n
}

func GetPossiblePieces(position Position, color int, type_ int) []Position {
	possiblePieces := []Position{}
	for r, row := range Board {
		for c, possiblePiece := range row {
			if possiblePiece.Color == color && possiblePiece.Type == type_ &&
				IsValidMove(Move{Position{r, c}, position}) {
				possiblePieces = append(possiblePieces, Position{r, c})
			}
		}
	}

	return possiblePieces
}
