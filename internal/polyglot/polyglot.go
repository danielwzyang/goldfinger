package polyglot

import (
	"bufio"
	"embed"
	"encoding/binary"
	"fmt"
	"math/rand/v2"

	"danielyang.cc/chess/internal/board"
)

const (
	BLACK_PAWN = iota
	WHITE_PAWN
	BLACK_KNIGHT
	WHITE_KNIGHT
	BLACK_BISHOP
	WHITE_BISHOP
	BLACK_ROOK
	WHITE_ROOK
	BLACK_QUEEN
	WHITE_QUEEN
	BLACK_KING
	WHITE_KING

	ENTRY_SIZE = 16
)

type Entry struct {
	Move   string
	Weight uint16
}

type MoveBin struct {
	Move    string
	Minimum int
	Maximum int
}

var (
	Book             = map[uint64][]Entry{}
	fileStrings      = []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	promotionStrings = []string{"", "n", "b", "r", "q"}

	//go:embed books/*.bin
	embeddedAssets embed.FS
)

func LoadBook(book string) {
	file, err := embeddedAssets.Open(book)

	if err != nil {
		panic("Cannot open file: " + err.Error())
	}

	buffer := make([]byte, ENTRY_SIZE)
	reader := bufio.NewReader(file)

	_, err = reader.Read(buffer)

	for err == nil {
		key, entry := parseEntry(buffer)

		Book[key] = append(Book[key], entry)

		_, err = reader.Read(buffer)
	}
}

func HasBookMove() bool {
	hash := GetPolyglotHash()
	entries, ok := Book[hash]

	if !ok {
		return false
	}

	return len(entries) > 0
}

func encodeMove(move string) int {
	return board.StringToMove(move)
}

func GetBestMove() int {
	hash := GetPolyglotHash()
	entries := Book[hash]

	return encodeMove(entries[0].Move)
}

func GetWeightedRandomMove() int {
	hash := GetPolyglotHash()
	entries := Book[hash]

	moveBins := make([]MoveBin, len(entries))

	acc := 0
	for i, entry := range entries {
		moveBins[i] = MoveBin{
			entry.Move,
			acc,
			acc + int(entry.Weight),
		}

		acc += int(entry.Weight)
	}

	random := rand.IntN(int(acc))

	for _, moveBin := range moveBins {
		if random >= moveBin.Minimum && random < moveBin.Maximum {
			return encodeMove(moveBin.Move)
		}
	}

	return encodeMove(entries[0].Move)
}

func parseEntry(bytes []byte) (uint64, Entry) {
	key := binary.BigEndian.Uint64(bytes[:8])
	move := binary.BigEndian.Uint16(bytes[8:10])
	weight := binary.BigEndian.Uint16(bytes[10:12])

	return key, Entry{parseMove(move), weight}
}

func parseMove(move uint16) string {
	promotion := move >> 12 & 7
	fromRow := move >> 9 & 7
	fromFile := move >> 6 & 7
	toRow := move >> 3 & 7
	toFile := move & 7

	moveString := fmt.Sprintf("%s%d%s%d%s", fileStrings[fromFile], fromRow+1, fileStrings[toFile], toRow+1, promotionStrings[promotion])

	return convertCastles(moveString)
}

func convertCastles(move string) string {
	switch move {
	case "e1h1":
		if board.LS1B(board.Bitboards[board.WHITE_KING]) == board.E1 {
			return "e1g1"
		}
	case "e1a1":
		if board.LS1B(board.Bitboards[board.WHITE_KING]) == board.E1 {
			return "e1c1"
		}
	case "e8h8":
		if board.LS1B(board.Bitboards[board.BLACK_KING]) == board.E8 {
			return "e8g8"
		}
	case "e8a8":
		if board.LS1B(board.Bitboards[board.BLACK_KING]) == board.E8 {
			return "e8c8"
		}
	}

	return move
}
