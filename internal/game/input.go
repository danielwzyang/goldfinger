package game

import (
	"fmt"
	"regexp"

	"danielyang.cc/chess/internal/board"
)

func alphaToNumeric(position string) [2]int {
	// converts string to rune array
	runes := []rune(position)

	// ascii value of 'a' is 97, '8' is 56
	return [2]int{56 - int(runes[1]), int(runes[0]) - 97}
}

func InputMove(color byte) {
	for {
		move := Input()

		// castling
		if move == "0-0" {
			if !board.ValidKingSideCastle(color) {
				fmt.Println("You cannot castle kingside.")
				continue
			}
			board.KingsideCastle(color)
			return
		}

		if move == "0-0-0" {
			if !board.ValidQueenSideCastle(color) {
				fmt.Println("You cannot castle queenside.")
				continue
			}
			board.QueensideCastle(color)
			return
		}

		// basic notation for moves
		if !fitsRegexPattern(move, "^([NBRQK]?)([a-h][1-8])$") {
			fmt.Println("Please follow the format described in /docs/input.md.")
			continue
		}

		var piece byte
		var finalPos [2]int
		promotion := false
		// moving pawn
		if len(move) == 2 {
			piece = 'P'
			finalPos = alphaToNumeric(move)

			if (color == 'w' && finalPos[0] == 0) || (color == 'b' && finalPos[0] == 7) {
				promotion = true
			}
		} else {
			piece = move[0]
			finalPos = alphaToNumeric(move[1:])
		}

		possiblePieces := board.GetPossiblePieces(color, piece, finalPos[0], finalPos[1])

		if len(possiblePieces) == 0 {
			fmt.Println("You own no such piece to make this move.")
			continue
		}

		if len(possiblePieces) == 2 {
			fmt.Println("Multiple pieces can make this move. Type the position of the piece you want to move.")
			move = Input()
			for !fitsRegexPattern(move, "^[a-h][1-8]$") || !board.ContainsPosition(possiblePieces, alphaToNumeric(move)) {
				fmt.Println("Please type a valid position.")
				move = Input()
			}
			numeric := alphaToNumeric(move)
			board.MakeMove(numeric[0], numeric[1], finalPos[0], finalPos[1], false)

			if promotion {
				fmt.Println("What piece do you want to promote your pawn to? (N | B | R | Q)")
				newPiece := Input()
				for !fitsRegexPattern(newPiece, "^[NBRQ]$") {
					fmt.Println("Please type a valid piece.")
					newPiece = Input()
				}
				board.Board[finalPos[0]][finalPos[1]] = string(color) + newPiece
			}

			return
		}

		board.MakeMove(possiblePieces[0][0], possiblePieces[0][1], finalPos[0], finalPos[1], false)

		if promotion {
			fmt.Println("What piece do you want to promote your pawn to? (N | B | R | Q)")
			newPiece := Input()
			for !fitsRegexPattern(newPiece, "^[NBRQ]$") {
				fmt.Println("Please type a valid piece. You can promote to N, B, R, or Q.")
				newPiece = Input()
			}
			board.Board[finalPos[0]][finalPos[1]] = string(color) + newPiece
		}

		return
	}
}

func fitsRegexPattern(str string, regexPattern string) bool {
	return regexp.MustCompile(regexPattern).MatchString(str)
}
