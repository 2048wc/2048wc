package boardLib_test

import "testing"
import "github.com/2048wc/2048wc/boardLib"

// test that the board is an n by n iterable where all elements are 0
func TestBoardInitialised(t *testing.T) {
	board := boardLib.CreateBoard()
	for i := 0; i < boardLib.BoardSize; i++ {
		for j := 0; j < boardLib.BoardSize; j++ {
			if board[i][j] != 0 {
				t.Fail()
			}
		}
	}
}

// we will mostly likely only use board size of 4. If one day
// we decide otherwise, we can update the test
func TestBoardSize(t *testing.T) {
	if boardLib.BoardSize != 4 {
		t.Fail()
	}
}
