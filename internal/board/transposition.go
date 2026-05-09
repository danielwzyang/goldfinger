package board

type NodeType uint8

const (
    PVNode NodeType = iota
    AllNode
    CutNode
    NoneNode
)

type TTEntry struct {
    Key   uint64
    Move  int
    Score int    
    Depth int     
    Type  NodeType 
}

const TT_SIZE = 1 << 20 
const TT_MASK = TT_SIZE - 1

var TRANSPOSITION_TABLE [TT_SIZE]TTEntry

func ResetTT() {
    TRANSPOSITION_TABLE = [TT_SIZE]TTEntry{}
}

func AddTTEntry(hash uint64, move int, score int, depth int, ply int, nodeType NodeType) {
	// board.MATE is 30000
    if score > 29000 {
        score += ply
    } else if score < -29000 {
        score -= ply
    }

	// index with mask is faster than mod
    index := hash & TT_MASK
    entry := &TRANSPOSITION_TABLE[index]
    
	// only overwrite if new search is deeper
    if entry.Key == 0 || depth >= entry.Depth {
        entry.Key = hash
        entry.Move = move
        entry.Score = score
        entry.Depth = depth
        entry.Type = nodeType
    }
}

func GetTTEntry(hash uint64, ply int) (TTEntry, bool) {
    index := hash & TT_MASK
    entry := TRANSPOSITION_TABLE[index]

    if entry.Key == hash {
        score := entry.Score
        if score > 29000 {
            score -= ply
        } else if score < -29000 {
            score += ply
        }
        entry.Score = score
        return entry, true
    }

    return TTEntry{}, false
}
