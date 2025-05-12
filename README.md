# Goldfinger

Based on games played against Komodo, this engine is probably rated around 1600-1700 given a search depth of 8-9 plies.

https://www.chess.com/member/goldfinger-1964

## Current Functionalities
- FEN Processing
- Pseudo Legal Move Generation with Magic Bitboards
- Alpha-Beta Pruning with the Negamax framework
- Late Move Reduction
- Null Move Pruning
- Move Ordering with MVV-LVA, Killer Heuristic, and History Heuristic
- PeSTO Evaluation
- Quiescence Search with Delta Pruning and basic Static Exchange Evaluation
- Tranposition Table with Zobrist Hashing
- Perft Test

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