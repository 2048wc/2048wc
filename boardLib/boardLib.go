package boardLib

import "fmt"

// The size of the board. Most likely will always stay 4,
// as this is what seems to be most playable.
const BoardSize = 4

// Creates an empty board of board_size. We have decided
// to use an array, becasue the size is immutable anyway.
func CreateBoard() (newBoard [BoardSize][BoardSize]int) {
	return
}

// This might be useful for logging/ command line interface.
// For the database we'll have to use a different function that's
// going to return json.
func PrintBoard(board [BoardSize][BoardSize]int) {
	for _, row := range board {
		for _, value := range row {
			fmt.Print(value, " ")
		}
		fmt.Println("")
	}
}
