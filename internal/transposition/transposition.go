package transposition

import "danielyang.cc/chess/internal/board"

type NodeType int

const (
	PVNode NodeType = iota
	AllNode
	CutNode
)

type Node struct {
	BestMove    board.Move
	Score       float64
	DepthLeft   int
	Type        NodeType
	SortedMoves []board.Move
}

var (
	table = map[uint64]Node{}
)

func Init() {
	initZobrist()
}

func AddEntry(nodeType NodeType, bestMove board.Move, score float64, depthLeft int, sortedMoves []board.Move, color int) {
	table[HashBoard(color)] = Node{
		bestMove,
		score,
		depthLeft,
		nodeType,
		sortedMoves,
	}
}

func GetEntry(color int) (Node, bool) {
	val, ok := table[HashBoard(color)]
	return val, ok
}
