package board

type NodeType int

const (
	PVNode NodeType = iota
	AllNode
	CutNode
)

type Node struct {
	Move  int
	Score int
	Depth int
	Type  NodeType
}

var TRANSPOSITION_TABLE = map[uint64]Node{}

func AddTTEntry(move int, score int, depth int, nodeType NodeType) {
	TRANSPOSITION_TABLE[ZobristHash] = Node{
		move,
		score,
		depth,
		nodeType,
	}
}

func GetTTEntry() (Node, bool) {
	val, ok := TRANSPOSITION_TABLE[ZobristHash]
	return val, ok
}
