package transposition

type NodeType int

const (
	PVNode NodeType = iota
	AllNode
	CutNode
)

type Node struct {
	BestMove  [2][2]int
	Score     float64
	DepthLeft int
	Type      NodeType
}

var (
	table = map[uint64]Node{}
)

func Init() {
	initZobrist()
}

func AddEntry(nodeType NodeType, bestMove [2][2]int, score float64, depthLeft int, color byte) {
	table[HashBoard(color)] = Node{Type: nodeType, BestMove: bestMove, Score: score, DepthLeft: depthLeft}
}

func GetEntry(color byte) (Node, bool) {
	val, ok := table[HashBoard(color)]
	return val, ok
}
