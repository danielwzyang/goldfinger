# Goldfinger

Based on games played against Komodo, this engine is probably rated around 1600-1700 given a search depth of 8-9 plies.

https://www.chess.com/member/goldfinger-1964

## Current Functionalities
- FEN Processing
- GUI for Board Representation
- Input Moves using Long Algebraic Notation
- Pseudo Legal Move Generation with Magic Bitboards
- Iterative Deepening with Aspiration Windows
- Alpha-Beta Pruning with the Negamax framework
- Late Move Reduction
- Null Move Pruning
- Move Ordering with MVV-LVA, Killer Heuristic, and History Heuristic
- PeSTO Evaluation
- Quiescence Search with Delta Pruning and basic Static Exchange Evaluation
- Tranposition Table with Zobrist Hashing
- Perft Test
- Repetition Table (disregards 3 in a row rule; just ignores the first repetition)
- Fifty Move Rule Detection

## Possible Improvements
- UCI Support / Link with Lichess API
- Opening / Ending Books
- NNUE Evaluation
- Material Draw
- More Pruning / Reduction Strategies e.g. Recursive Static Exchange Evaluation

## Usage

Ensure Go is installed in the latest release.

To play against then engine, run cmd/goldfinger/main.go using the following optional flags to customize behavior/performance:

|Flag|Description|Default|Usage|
|-|-|-|-|
|fen|board position|starting position|-fen="rnbqkbnr/pp..."
|depth|search depth|6|-depth=8|
|black|if included input controls black|no flag = input controls white|-black|

Input moves using long algebraic notation.

To watch the engine play against itself, simply run cmd/self/main.go and use any of the optional flags, disregarding -black. 

## Games Against Komodo Engine

|Game #|Opponent|Result|Accuracy|Depth|Estimated ELO|Notes|
|-|-|-|-|-|-|-|
|1|Komodo1 (250)|Win|88.8%|6|1600|
|2|Komodo2 (400)|Win|77.9%|6|1400|
|3|Komodo3 (550)|Win|87.0%|6|1700|
|4|Komodo4 (700)|Win|80.2%|6|1600|
|5|Komodo5 (850)|Win|91.8%|6|1900|
|6|Komodo6 (1000)|Win|81.9%|6|1750|
|7|Komodo7 (1100)|Draw by repetition|87.6%|8|1950|
|8|Komodo7 (1100)|Win|89.5%|8|2000|
|9|Komodo8 (1200)|Win|83.2%|8|1900|
|10|Komodo9 (1300)|Win|88.5%|8|2050|
|11|Komodo10 (1400)|Win|79.4%|8|1750|
|12|Komodo11 (1500)|Win|89.7%|8|2200|
|13|Komodo12 (1600)|Win|78.6%|8|1750|
|14|Komodo13 (1700)|Loss|77.7%|8|1800|
|15|Komodo13 (1700)|Loss|70.8%|8|1350|
|16|Komodo13 (1700)|Loss|74.0%|8|1550|
|17|Komodo13 (1700)|Draw by repetition|75.7%|9|1600|
|18|Komodo13 (1700)|Win|68.8%|9|1250|testing changes; bad performance
|19|Komodo13 (1700)|Win|83.6%|9|2050|
|20|Komodo14 (1800)|Win|76.5%|9|1700|
|21|Komodo15 (1900)|Loss|69.4%|9|1100|
|22|Komodo14 (1800)|Loss|61.4%|8|1050|worst endgame by far (45.7%); not really sure what happened
|23|Komodo13 (1700)|Loss|51.5%|8|950|very bad game + first game on white; removing iterative deepening after this
|24|Komodo13 (1700)|Loss|63.8%|8|1100|bad game, also on white; no iterative deepening
|25|Komodo13 (1700)|Win|81.5%|8|1950|after figuring out major flaw in evaluation; readded iterative deepening
|26|Komodo14 (1800)|Win|88.5%|9|2300|longest thinking time of every game so far: (Avg: 6311ms \| Max: 117253ms \| Total: 258775ms); also first brilliant move!!