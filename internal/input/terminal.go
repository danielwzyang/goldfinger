package input

import (
	"fmt"
)

func Clear() {
	fmt.Print("\033[H\033[J")
	header()
}

func Input() string {
	fmt.Printf("> ")

	var input string
	_, err := fmt.Scanln(&input)

	for err != nil {
		fmt.Scanln(&input)
	}

	return input
}

func header() {
	fmt.Println("──────────────────────────────────────────────────────")
	fmt.Println("Pawn Star | danielyang.cc ©")
	fmt.Println("──────────────────────────────────────────────────────")
	fmt.Println()
}

func StartOpts(player *byte) {
	Clear()

	fmt.Println("What color are you playing? (b | w)")
	fmt.Println("The opposing team will be played by the engine/bot.")
	color := Input()
	for color != "b" && color != "w" {
		fmt.Println("Please enter b or w.")
		color = Input()
	}
	*player = color[0]
}
