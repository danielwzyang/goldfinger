# Goldfinger

Rated ~2000 on Lichess

https://lichess.org/@/goldfinger-bot

## Current Functionalities
- UCI Support + Linked with Lichess
- FEN Processing
- Input Moves using Long Algebraic Notation
- Move Encoding (int32)
- GUI for Board Representation
- Pseudo Legal Move Generation with Magic Bitboards
- Iterative Deepening with Aspiration Windows
- Alpha-Beta Pruning with the Negamax framework
- Principal Variation Search
- Late Move Reduction
- Null Move Pruning
- Move Ordering with MVV-LVA, Killer Heuristic, and History Heuristic
- PeSTO Evaluation
- Quiescence Search with Delta Pruning and Static Exchange Evaluation
- Tranposition Table with Zobrist Hashing
- Perft Test
- Repetition Table
- Fifty Move Rule Detection
- Polyglot Opening Book Support

## Possible Improvements
- NNUE Evaluation
- Better Handcrafted Evaluation with Pawn Structure, King Safety, etc.

## Usage

Ensure Go is installed in the latest release.

To use the UCI, run cmd/uci/main.go and use supported UCI commands.

To play against the engine, run cmd/goldfinger/main.go using the following optional flags to customize behavior/performance:

|Flag|Description|Default|Usage|
|-|-|-|-|
|fen|board position|starting position|-fen="rnbqkbnr/pp..."
|time|search time in ms|1000|-time=100|
|black|if included input controls black|no flag = input controls white|-black|

Input moves using long algebraic notation.

To watch the engine play against itself, simply run cmd/self/main.go and use any of the optional flags, disregarding -black. 
