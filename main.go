package main

import "github.com/2048wc/2048wc/boardLib"

func main() {
	board := boardLib.CreateBoard()
	board[2][3] = 4
	board[0][1] = 2
	boardLib.PrintBoard(board)
}
