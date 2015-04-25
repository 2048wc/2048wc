/* Copyright (C) 2015  Adam Kurkiewicz and Iva Babukova
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published
 *   by the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package boardLib

import "testing"
import "encoding/json"
import "log"

// test that the board is an n by n iterable where all elements are 0
func TestBoardInitialised(t *testing.T) {
	var board BoardT
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			if board[i][j] != 0 {
				t.Fail()
			}
		}
	}
}

// we will mostly likely only use board size of 4. If one day
// we decide otherwise, we can update the test
func TestBoardSize(t *testing.T) {
	if BoardSize != 4 {
		t.Fail()
	}
}

func setUpMoveRightOldBoard() (oldBoard BoardT) {
	oldBoard = [BoardSize][BoardSize]int{
		{16, 8, 4, 2},
		{4, 2, 2, 0},
		{2, 4, 0, 2},
		{0, 0, 0, 0},
	}
	return
}

func setUpMoveRightNewBoard() BoardT {
	return [BoardSize][BoardSize]int{
		{16, 8, 4, 2},
		{0, 0, 4, 4},
		{0, 2, 4, 2},
		{0, 0, 0, 0},
	}
}

func setUpMoveRightNonMergedMoves() []moveT {
	return []moveT{
		moveT{{1, 0}, {1, 2}},
		moveT{{2, 0}, {2, 1}},
		moveT{{2, 1}, {2, 2}},
	}
}

func setUpMoveRightNewTileCandidates() []positionT {
	return []positionT{{1, 0}, {2, 0}, {3, 0}, {3, 1}, {3, 2}, {3, 3}}
}

func setUpMoveRightNonMovedTiles() []positionT {
	return []positionT{{0, 0}, {0, 1}, {0, 2}, {0, 3}}
}

func setUpMoveRightMergeMoves() []moveValueT {
	return []moveValueT{{moveT{{1, 1}, {1, 2}}, 4}}
}

func setUpMoveRightRandomTile() (randomTile []positionValueT) {
	return []positionValueT{{positionT{1, 1}, 2}}
}

func setUpMoveRightDirection() string {
	return "right"
}

func setUpMoveRightSeed() string {
	return "e9ccc20fdb924ed423ad1b46c6df43516685f4c2bc36e202ad467af1b1d1febf"
}

func setUpMoveRightRoundNo() int {
	return 24
}

func setUpMoveIsGameOver() bool {
	return false
}

func setUpRightMoveHarness() Move {
	return Move{
		Direction:         setUpMoveRightDirection(),
		Seed:              setUpMoveRightSeed(),
		IsGameOver:        setUpMoveIsGameOver(),
		RoundNo:           setUpMoveRightRoundNo() + 1,
		OldBoard:          setUpMoveRightOldBoard(),
		NewTileCandidates: setUpMoveRightNewTileCandidates(),
		NonMergeMoves:     setUpMoveRightNonMergedMoves(),
		MergeMoves:        setUpMoveRightMergeMoves(),
		NonMovedTiles:     setUpMoveRightNonMovedTiles(),
		NewBoard:          setUpMoveRightNewBoard(),
		RandomTiles:       setUpMoveRightRandomTile(),
	}

}

// This test is a simple assertion to watch the Move struct, which a lot of
// infrastructure (database, client) rely on. If the struct changes, or more
// precisely, if struct jsonification changes, then this test should fail
func TestJsonMarshalling(t *testing.T) {
	move := setUpRightMoveHarness()
	marshalled, err := json.Marshal(move)
	if err != nil {
		log.Fatal(err)
	}
	if string(marshalled) != `{"Direction":"right","RoundNo":25,"Seed":"e9ccc20fdb924ed423ad1b46c6df43516685f4c2bc36e202ad467af1b1d1febf","OldBoard":[[16,8,4,2],[4,2,2,0],[2,4,0,2],[0,0,0,0]],"NewBoard":[[16,8,4,2],[0,0,4,4],[0,2,4,2],[0,0,0,0]],"NonMergeMoves":[[[1,0],[1,2]],[[2,0],[2,1]],[[2,1],[2,2]]],"MergeMoves":[{"Move":[[1,1],[1,2]],"Value":4}],"NonMovedTiles":[[0,0],[0,1],[0,2],[0,3]],"NewTileCandidates":[[1,0],[2,0],[3,0],[3,1],[3,2],[3,3]],"RandomTiles":[{"Position":[1,1],"Value":2}],"IsGameOver":false}` {
		t.Fail()
	}
}

func TestIterator(t *testing.T) {
	board := setUpMoveRightOldBoard()
	iter := boardIterator(board, "right")
	var boardIndex boardIndexT
	var err error
	expectedCurrentIndex := [][2]int{{0, 3}, {0, 2}, {0, 1}, {0, 0}, {1, 3},
		{1, 2}, {1, 1}, {1, 0}, {2, 3}, {2, 2}, {2, 1}, {2, 0}, {3, 3}, {3, 2},
		{3, 1}}
	expectedNextIndex := [][2]int{{0, 2}, {0, 1}, {0, 0}, {1, 3},
		{1, 2}, {1, 1}, {1, 0}, {2, 3}, {2, 2}, {2, 1}, {2, 0}, {3, 3}, {3, 2},
		{3, 1}, {3, 0}}
	for i := 0; i < len(expectedCurrentIndex); i++ {
		boardIndex, err = iter()
		if err != nil {
			t.Error(err)
		}
		if expectedCurrentIndex[i] != boardIndex.currentIndex {
			t.Error("error at current index", i, expectedCurrentIndex[i])
		}
		if expectedNextIndex[i] != boardIndex.nextIndex {
			t.Error("error at next index", i, expectedCurrentIndex[i])
		}
	}
}

// This test the first stage of the pipeline
func TestComputeDistance(t *testing.T) {
	move := CreateMove(
		setUpMoveRightOldBoard(),
		setUpMoveRightDirection(),
		setUpMoveRightRoundNo(),
		setUpMoveRightSeed(),
	)
	expected := [4][4]int{
		{0, 0, 0, 0},
		{2, 2, 1, 0},
		{1, 1, 0, 0},
		{3, 2, 1, 0},
	}
	distances, err := move.computeDistance()
	if distances != expected {
		t.Error(expected, "!=", distances)
	}
	if err != nil {
		t.Error(err, "is not nil")
	}

}

// This test exposes a bug, which slipped through to master.
func TestComputeDistanceLastRow(t *testing.T) {
	move := CreateMove(
		[BoardSize][BoardSize]int{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{2, 0, 0, 0},
		},
		setUpMoveRightDirection(),
		setUpMoveRightRoundNo(),
		setUpMoveRightSeed(),
	)
	distances, err := move.computeDistance()
	if distances[0][0] != 3 {
		t.Error("distances should be 3")
	}
	if err != nil {
		t.Error(err, "is not nil")
	}

}

// This function is testing the Move pipeline. It needs to make sure that the
// pipeline modifies a move struct from a given, non-specific point
// appropriately.
func TestMoveRight(t *testing.T) {
	return
	move := CreateMove(
		setUpMoveRightOldBoard(),
		setUpMoveRightDirection(),
		setUpMoveRightRoundNo(),
		setUpMoveRightSeed(),
	)
	move.ExecuteMove()
	result, err := json.Marshal(move)
	if err != nil {
		log.Fatal(err)
	}
	moveHarness := setUpRightMoveHarness()
	harness, err := json.Marshal(moveHarness)
	if err != nil {
		log.Fatal(err)
	}
	if string(result) != string(harness) {
		t.Fail()
	}
}
