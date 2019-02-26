package main

import (
	"fmt"

	"github.com/crepehat/fourplay/fourplay"
)

// 14363756335665245414
func main() {
	board, err := fourplay.CreateFromSequence("14363756335665245414444")
	if err != nil {
		fmt.Println(board)
	}
	// board.Print()
}
