package board

type MoveList struct {
	Moves [256]int
	Count int
}

func (moveList *MoveList) AddMove(move int) {
	moveList.Moves[moveList.Count] = move
	moveList.Count++
}

func EncodeMove(source int, target int, piece int, promotion int, capture int, double int, enpassant int, castling int) int {
	return (source) | // 6 bits 2^6 = 64
		(target << 6) | // 6 bits
		(piece << 12) | // 4 bits 2^4 = 16 (12 pieces)
		(promotion << 16) | // 4 bits (12 pieces)
		(capture << 20) | // 1 bit (bool)
		(double << 21) | // 1 bit
		(enpassant << 22) | // 1 bit
		(castling << 23) // 1 bit
}

func GetSource(move int) int {
	return move & 0x3f
}

func GetTarget(move int) int {
	return (move & 0xfc0) >> 6
}

func GetPiece(move int) int {
	return (move & 0xf000) >> 12
}

func GetPromotion(move int) int {
	return (move & 0xf0000) >> 16
}

func GetCapture(move int) int {
	return move & 0x100000
}

func GetDouble(move int) int {
	return move & 0x200000
}

func GetEnPassant(move int) int {
	return move & 0x400000
}

func GetCastling(move int) int {
	return move & 0x800000
}
