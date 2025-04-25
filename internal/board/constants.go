package board

type BoardState struct {
	LastMove Move
	PieceA   Piece // piece that was moved
	PieceB   Piece // piece that was "captured" / empty piece

	WCastleKS bool
	WCastleQS bool
	BCastleKS bool
	BCastleQS bool

	EnPassant Position
}

type Position struct {
	Rank int
	File int
}

type Move struct {
	From Position
	To   Position
}

type Piece struct {
	Type  int
	Color int
	Key   int
}

func EMPTY_PIECE() Piece {
	return Piece{EMPTY, EMPTY, 0}
}

func WHITE_PAWN() Piece {
	return Piece{PAWN, WHITE, 1}
}

func WHITE_KNIGHT() Piece {
	return Piece{KNIGHT, WHITE, 2}
}

func WHITE_BISHOP() Piece {
	return Piece{BISHOP, WHITE, 3}
}

func WHITE_ROOK() Piece {
	return Piece{ROOK, WHITE, 4}
}

func WHITE_QUEEN() Piece {
	return Piece{QUEEN, WHITE, 5}
}

func WHITE_KING() Piece {
	return Piece{KING, WHITE, 6}
}

func BLACK_PAWN() Piece {
	return Piece{PAWN, BLACK, 7}
}

func BLACK_KNIGHT() Piece {
	return Piece{KNIGHT, BLACK, 8}
}

func BLACK_BISHOP() Piece {
	return Piece{BISHOP, BLACK, 9}
}

func BLACK_ROOK() Piece {
	return Piece{ROOK, BLACK, 10}
}

func BLACK_QUEEN() Piece {
	return Piece{QUEEN, BLACK, 11}
}

func BLACK_KING() Piece {
	return Piece{KING, BLACK, 12}
}

var (
	DefaultBoard = [8][8]Piece{
		{BLACK_ROOK(), BLACK_KNIGHT(), BLACK_BISHOP(), BLACK_QUEEN(), BLACK_KING(), BLACK_BISHOP(), BLACK_KNIGHT(), BLACK_ROOK()},
		{BLACK_PAWN(), BLACK_PAWN(), BLACK_PAWN(), BLACK_PAWN(), BLACK_PAWN(), BLACK_PAWN(), BLACK_PAWN(), BLACK_PAWN()},
		{EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE()},
		{EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE()},
		{EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE()},
		{EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE(), EMPTY_PIECE()},
		{WHITE_PAWN(), WHITE_PAWN(), WHITE_PAWN(), WHITE_PAWN(), WHITE_PAWN(), WHITE_PAWN(), WHITE_PAWN(), WHITE_PAWN()},
		{WHITE_ROOK(), WHITE_KNIGHT(), WHITE_BISHOP(), WHITE_QUEEN(), WHITE_KING(), WHITE_BISHOP(), WHITE_KNIGHT(), WHITE_ROOK()},
	}

	//				  0    1     2     3     4    5     6     7     8     9    10    11    12
	ascii = []string{" ", "♙", "♘", "♗", "♖", "♕", "♔", "♟", "♞", "♝", "♜", "♛", "♚"}

	EMPTY = -1

	PAWN          = 0
	KNIGHT        = 1
	BISHOP        = 2
	ROOK          = 3
	QUEEN         = 4
	KING          = 5
	PIECE_VALUES  = map[byte]int{'P': PAWN, 'N': KNIGHT, 'B': BISHOP, 'R': ROOK, 'Q': QUEEN, 'K': KING}
	PIECE_LETTERS = [6]string{"P", "N", "B", "R", "Q", "K"}

	WHITE = 0
	BLACK = 1

	MATE_SCORE  = 100000
	LIMIT_SCORE = 1000000
)

func GetKey(piece int, color int) int {
	return color*6 + piece + 1
}
