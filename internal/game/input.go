package game

import (
	"fmt"
	"regexp"

	"danielyang.cc/chess/internal/board"
)

func alphaToNumeric(position string) board.Position {
	// converts string to rune array
	runes := []rune(position)

	// ascii value of 'a' is 97, '8' is 56
	return board.Position{Rank: 56 - int(runes[1]), File: int(runes[0]) - 97}
}

func InputMove() {
	for {
		move := Input()

		// castling
		if move == "0-0" {
			if !board.ValidKingSideCastle(playerColor) {
				fmt.Println("You cannot castle kingside.")
				continue
			}
			board.KingsideCastle(playerColor)
			return
		}

		if move == "0-0-0" {
			if !board.ValidQueenSideCastle(playerColor) {
				fmt.Println("You cannot castle queenside.")
				continue
			}
			board.QueensideCastle(playerColor)
			return
		}

		// basic notation for moves
		if !fitsRegexPattern(move, "^([NBRQK]?)([a-h][1-8])$") {
			fmt.Println("Please follow the format described in /docs/input.md.")
			continue
		}

		var piece int
		var to board.Position
		promotion := false
		// moving pawn
		if len(move) == 2 {
			piece = board.PAWN

			to = alphaToNumeric(move)

			// promoting pawn
			if (playerColor == board.WHITE && to.Rank == 0) || (playerColor == board.BLACK && to.Rank == 7) {
				promotion = true
			}
		} else {
			piece = board.PIECE_VALUES[move[0]]
			to = alphaToNumeric(move[1:])
		}

		// getting pieces that could move to the described position
		possiblePieces := board.GetPossiblePieces(to, playerColor, piece)

		if len(possiblePieces) == 0 {
			fmt.Println("You cannot make this move.")
			continue
		}

		if len(possiblePieces) > 1 {
			fmt.Println("Multiple pieces can make this move. Type the position of the piece you want to move.")

			// input position for disambiguation
			move = Input()
			for !fitsRegexPattern(move, "^[a-h][1-8]$") || !board.ContainsPosition(possiblePieces, alphaToNumeric(move)) {
				fmt.Println("Please type a valid position.")
				move = Input()
			}
			from := alphaToNumeric(move)
			board.MakeMove(board.Move{From: from, To: to})

			// handle promotion
			if promotion {
				// input piece
				fmt.Println("What piece do you want to promote your pawn to? (N | B | R | Q)")
				newPiece := Input()
				for !fitsRegexPattern(newPiece, "^[NBRQ]$") {
					fmt.Println("Please type a valid piece.")
					newPiece = Input()
				}

				piece = board.PIECE_VALUES[newPiece[0]]

				// update piece
				board.Board[to.Rank][to.File] = board.Piece{
					Type:  piece,
					Color: playerColor,
					Key:   board.GetKey(piece, playerColor),
				}
			}

			return
		}

		// there's only one possible piece that can make the move
		board.MakeMove(board.Move{From: possiblePieces[0], To: to})

		// handle promotion
		if promotion {
			// input piece
			fmt.Println("What piece do you want to promote your pawn to? (N | B | R | Q)")
			newPiece := Input()
			for !fitsRegexPattern(newPiece, "^[NBRQ]$") {
				fmt.Println("Please type a valid piece. You can promote to N, B, R, or Q.")
				newPiece = Input()
			}

			piece = board.PIECE_VALUES[newPiece[0]]

			// update piece
			board.Board[to.Rank][to.File] = board.Piece{
				Type:  piece,
				Color: playerColor,
				Key:   board.GetKey(piece, playerColor),
			}
		}

		return
	}
}

func fitsRegexPattern(str string, regexPattern string) bool {
	return regexp.MustCompile(regexPattern).MatchString(str)
}
