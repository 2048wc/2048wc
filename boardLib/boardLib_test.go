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
func testBoardInitialised(t *testing.T) {
	var board boardT
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			if board[i][j] != 0 {
				t.Fail()
			}
		}
	}
}

type moveTest struct{
	harness moveT
	result moveT
}
// we will mostly likely only use board size of 4. If one day
// we decide otherwise, we can update the test
func testBoardSize(t *testing.T) {
	if BoardSize != 4 {
		t.Fail()
	}
}

func setUpMoveRightOldBoard() (oldBoard boardT) {
	oldBoard = [BoardSize][BoardSize]int{
		{16, 8, 4, 2},
		{4, 2, 2, 0},
		{2, 4, 0, 2},
		{2, 0, 0, 0},
	}
	return
}

func setUpMoveRightNewBoard() boardT {
	return [BoardSize][BoardSize]int{
		{16, 8, 4, 2},
		{0, 0, 4, 4},
		{0, 2, 4, 2},
		{0, 0, 0, 2},
	}
}

func setUpMoveRightNonMergedMoves() []nonMergeMoveT {
	return []nonMergeMoveT{
		nonMergeMoveT{{1, 0}, {1, 2}},
		nonMergeMoveT{{2, 1}, {2, 2}},
		nonMergeMoveT{{2, 0}, {2, 1}},
		nonMergeMoveT{{3, 0}, {3, 3}},
	}
}

func setUpMoveRightNewTileCandidates() []positionT {
	return []positionT{{1, 0}, {2, 0}, {3, 0}, {3, 1}, {3, 2}, {3, 3}}
}

func setUpMoveRightNonMovedTiles() []positionT {
	return []positionT{{0, 3},{0, 2},{0, 1},{0, 0},{2, 3}}
}

func setUpMoveRightMergeMoves() []mergeMoveT {
	return []mergeMoveT{mergeMoveT{nonMergeMoveT{{1, 2}, {1, 1}}, positionT{1, 3}, 4}}
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

func setUpRightMoveHarness() moveT {
	return moveT{
		Direction:         setUpMoveRightDirection(),
		Seed:              setUpMoveRightSeed(),
		IsGameOver:        setUpMoveIsGameOver(),
		RoundNo:           setUpMoveRightRoundNo() + 1,
		OldBoard:          setUpMoveRightOldBoard(),
		NewBoard: 		   setUpMoveRightNewBoard(),
		NonMergeMoves:     setUpMoveRightNonMergedMoves(),
		MergeMoves:        setUpMoveRightMergeMoves(),
		NonMovedTiles:     setUpMoveRightNonMovedTiles(),
		NewTileCandidates: make([]positionT, 0, BoardSize*BoardSize),
		//NewBoard:          setUpMoveRightNewBoard(),
		//RandomTiles:       setUpMoveRightRandomTile(),
		RandomTiles: make([]positionValueT, 0, BoardSize*BoardSize),

	}
}

// This test is a simple assertion to watch the Move struct, which a lot of
// infrastructure (database, client) rely on. If the struct changes, or more
// precisely, if struct jsonification changes, then this test should fail.
//TODO rewrite Marshalling to produce a unique json representation.
func testJsonMarshalling(t *testing.T) {
	println("TODO implement")
	t.Skip()
	move := setUpRightMoveHarness()
	marshalled, err := json.Marshal(move)
	if err != nil {
		log.Fatal(err)
	}
	expected := `{"Direction":"right","RoundNo":25,"Seed":"e9ccc20fdb924ed423ad1b46c6df43516685f4c2bc36e202ad467af1b1d1febf","OldBoard":[[16,8,4,2],[4,2,2,0],[2,4,0,2],[0,0,0,0]],"NewBoard":[[16,8,4,2],[0,0,4,4],[0,2,4,2],[0,0,0,0]],"NonMergeMoves":[[[1,0],[1,2]],[[2,0],[2,1]],[[2,1],[2,2]]],"MergeMoves":[{"Move":[[1,1],[1,2]],"Value":4}],"NonMovedTiles":[[0,0],[0,1],[0,2],[0,3]],"NewTileCandidates":[[1,0],[2,0],[3,0],[3,1],[3,2],[3,3]],"RandomTiles":[{"Position":[1,1],"Value":2}],"IsGameOver":false}`
	if string(marshalled) != expected {
		t.Error(string(marshalled) + "\n" + expected)
	}
}

func testPopulateNewBoard(t *testing.T) {
	var move moveT
	move.InitNextMove(
		setUpMoveRightOldBoard(),
		setUpMoveRightDirection(),
		setUpMoveRightRoundNo(),
		setUpMoveRightSeed(),
	)
	expected := setUpMoveRightNewBoard()
	move.ResolveMove()
	if (&move).NewBoard != expected {
		t.Error(expected, "!=", (&move).NewBoard)
	}
}

// This function is testing the Move pipeline. It needs to make sure that the
// pipeline modifies a move struct from a given, non-specific point
// appropriately.

type executeMoveUnitTest struct {
	forEvaluation Move
	harness Move
}

func TestMoveRight(t *testing.T) {
	var move moveT
	move.InitNextMove(
		[BoardSize][BoardSize]int{
		{16, 8, 4, 2},
		{4, 2, 2, 0},
		{2, 4, 0, 2},
		{2, 0, 0, 0},
	},
		"right",
		24,
		"e9ccc20fdb924ed423ad1b46c6df43516685f4c2bc36e202ad467af1b1d1febf",
	)
	move.ResolveMove()
	moveHarness := setUpRightMoveHarness()
	v1, _ := json.Marshal(moveHarness)
	v2, _ := json.Marshal(move)
	println(string(v1) == string(v2))
	if string(v1) != string(v2) {
		t.Error("\nresult:    ", string(v1), "\nharness:", string(v2))
	}
}
