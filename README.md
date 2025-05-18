# Goldfinger

Rated ~1600-1700 on Lichess

https://lichess.org/@/goldfinger-bot
https://www.chess.com/member/goldfinger-1964

## Current Functionalities
- UCI Support + Linked with Lichess
- FEN Processing
- GUI for Board Representation
- Input Moves using Long Algebraic Notation
- Pseudo Legal Move Generation with Magic Bitboards
- Iterative Deepening with Aspiration Windows
- Alpha-Beta Pruning with the Negamax framework
- Late Move Reduction
- Null Move Pruning
- Move Ordering with MVV-LVA, Killer Heuristic, and History Heuristic
- Tapered Evaluation with Piece Square Tables, Pawn Structure, and Tapered King Safety/Activation
- Quiescence Search with Delta Pruning and basic Static Exchange Evaluation
- Tranposition Table with Zobrist Hashing
- Perft Test
- Repetition Table (disregards 3 in a row rule; just ignores the first repetition)
- Fifty Move Rule Detection

## Possible Improvements
- Opening / Ending Books
- NNUE Evaluation
- Material Draw
- More Pruning / Reduction Strategies e.g. Recursive Static Exchange Evaluation

## Usage

Ensure Go is installed in the latest release.

To use the UCI, run cmd/uci/main.go and use supported UCI commands.

To play against the engine, run cmd/goldfinger/main.go using the following optional flags to customize behavior/performance:

|Flag|Description|Default|Usage|
|-|-|-|-|
|fen|board position|starting position|-fen="rnbqkbnr/pp..."
|depth|search depth|6|-depth=8|
|black|if included input controls black|no flag = input controls white|-black|

Input moves using long algebraic notation.

To watch the engine play against itself, simply run cmd/self/main.go and use any of the optional flags, disregarding -black. 
