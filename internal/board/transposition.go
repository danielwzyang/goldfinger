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

func AddTTEntry(move int, score int, depth int, ply int, nodeType NodeType) {
	// board.MATE is 30000
	if score > 29000 {
		score += ply
	} else if score < -29000 {
		score -= ply
	}
	TRANSPOSITION_TABLE[ZobristHash] = Node{
		move,
		score,
		depth,
		nodeType,
	}
}

func GetTTEntry(ply int) (Node, bool) {
	val, ok := TRANSPOSITION_TABLE[ZobristHash]
	if ok {
		if val.Score > 29000 {
			val.Score -= ply
		} else if val.Score < -29000 {
			val.Score += ply
		}
	}
	return val, ok
}
