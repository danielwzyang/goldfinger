package game

import (
	"fmt"
)

func Input() string {
	fmt.Printf("> ")

	var input string
	_, err := fmt.Scanln(&input)

	for err != nil {
		fmt.Scanln(&input)
	}

	return input
}

func Header() {
	fmt.Println()
	fmt.Println("──────────────────────────────────────────────────────")
	fmt.Println("Chess Engine | danielyang.cc ©")
	fmt.Println("──────────────────────────────────────────────────────")
	fmt.Println()
}

func StartOpts(player *byte, engineType *byte) {
	Header()

	fmt.Println("What color are you playing? (b | w)")
	fmt.Println("The opposing team will be played by the engine/bot.")
	color := Input()
	for color != "b" && color != "w" {
		fmt.Println("Please enter b or w.")
		color = Input()
	}
	*player = color[0]

	fmt.Println()

	fmt.Println("What type of engine do you want? (r | n)")
	fmt.Println("An 'r' type engine will play random moves. An 'n' type engine will employ an algorithm.")
	type_ := Input()
	for type_ != "r" && type_ != "n" {
		fmt.Println("Please enter r or n.")
		type_ = Input()
	}
	*engineType = type_[0]
}
