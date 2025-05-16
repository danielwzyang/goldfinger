package gui

import (
	"danielyang.cc/chess/internal/board"

	"embed"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Game struct {
	pieces    [12]*ebiten.Image
	bitboards [12]uint64
	lastMove  int
}

var (
	squareSize   = 60
	screenWidth  = squareSize * 8
	screenHeight = squareSize * 8

	lightSquare     = color.RGBA{235, 236, 208, 255}
	darkSquare      = color.RGBA{119, 149, 86, 255}
	highlightSource = color.RGBA{219, 203, 107, 255}
	highlightTarget = color.RGBA{219, 187, 112, 255}

	images = [12]string{
		"assets/white_pawn.png",
		"assets/white_knight.png",
		"assets/white_bishop.png",
		"assets/white_rook.png",
		"assets/white_queen.png",
		"assets/white_king.png",
		"assets/black_pawn.png",
		"assets/black_knight.png",
		"assets/black_bishop.png",
		"assets/black_rook.png",
		"assets/black_queen.png",
		"assets/black_king.png",
	}

	//go:embed assets/*.png
	embeddedAssets embed.FS

	// stores information about game
	gameInstance *Game
)

// start
func Run() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Goldfinger | danielyang.cc")

	err := ebiten.RunGame(gameInstance)
	if err != nil {
		panic(err)
	}
}

// update bitboards + last move
func UpdateBoard(move int) {
	gameInstance.bitboards = board.Bitboards
	gameInstance.lastMove = move
}

// create new instance and load embedded assets
func NewGame() {
	game := &Game{
		pieces:    [12]*ebiten.Image{},
		bitboards: board.Bitboards,
	}

	for piece, image := range images {
		file, err := embeddedAssets.Open(image)
		if err != nil {
			panic(err)
		}

		img, _, err := ebitenutil.NewImageFromReader(file)
		if err != nil {
			panic(err)
		}

		game.pieces[piece] = img
	}

	gameInstance = game
}

func (g *Game) Draw(screen *ebiten.Image) {
	source := board.GetSource(g.lastMove)
	target := board.GetTarget(g.lastMove)

	if g.lastMove == 0 {
		source = -1
		target = -1
	}

	var square int
	var light bool

	// draw colored squares
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square = (7-rank)*8 + file
			light = (rank+file)%2 == 0

			// determine square color
			var squareColor color.RGBA
			if square == source {
				squareColor = highlightSource
			} else if square == target {
				squareColor = highlightTarget
			} else if light {
				squareColor = lightSquare
			} else {
				squareColor = darkSquare
			}

			screen.SubImage(image.Rect(file*squareSize, rank*squareSize, (file+1)*squareSize, (rank+1)*squareSize)).(*ebiten.Image).Fill(squareColor)

			// file and rank labels
			var textColor color.RGBA
			if light {
				textColor = darkSquare
			} else {
				textColor = lightSquare
			}
			if file == 0 {
				rankLabel := string('1' + rune(7-rank))
				x := 5
				y := rank*squareSize + 13
				text.Draw(screen, rankLabel, basicfont.Face7x13, x, y, textColor)
			}
			if rank == 7 {
				fileLabel := string('a' + rune(file))
				x := (file+1)*squareSize - 10
				y := (rank+1)*squareSize - 3
				text.Draw(screen, fileLabel, basicfont.Face7x13, x, y, textColor)
			}

		}
	}

	// draw pieces
	var file int
	var rank int

	for piece, bitboard := range g.bitboards {
		for square := 0; square < 64; square++ {
			if board.GetBit(bitboard, square) == 1 {
				file = square % 8
				rank = 7 - (square / 8)

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(file*squareSize), float64(rank*squareSize))
				screen.DrawImage(g.pieces[piece], op)
			}
		}
	}
}

// required function
func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
