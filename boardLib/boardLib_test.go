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

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

// test that the board is an n by n iterable where all elements are 0
func TestBoardInitialised(t *testing.T) {
	var board boardT
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

type moveTest struct {
	task           moveT
	result         moveT
	expectedResult moveT
}

func (move *moveTest) check() bool {
	var result bool
	result = reflect.DeepEqual(move.result, move.expectedResult)
	if result == false {
		fmt.Fprint(os.Stderr,
			"\ntask:          ", move.task,
			"\nresult:        ", move.result,
			"\nexpectedResult:", move.expectedResult)
	}
	return result
}

//TODO test other fields
func firstPassAllDirectionsNewBoardOnly() []moveTest {
	finishOffMoveTest := func(test *moveTest) {
		fullResult := test.task
		fullResult.firstPass()
		test.result = test.task
		test.expectedResult.OldBoard = test.task.OldBoard
		test.expectedResult.Direction = test.task.Direction
		test.result.NewBoard = fullResult.NewBoard
	}
	tests := make([]moveTest, 0)
	moveTestRight := moveTest{
		task: moveT{
			OldBoard: [BoardSize][BoardSize]int{
				{2, 2, 2, 2},
				{4, 2, 2, 0},
				{2, 0, 2, 0},
				{2, 0, 0, 2},
			},
			Direction: "right",
		},
		expectedResult: moveT{
			NewBoard: [BoardSize][BoardSize]int{
				{0, 0, 4, 4},
				{0, 0, 4, 4},
				{0, 0, 0, 4},
				{0, 0, 0, 4},
			},
		},
	}
	moveTestDown := moveTest{
		task: moveT{
			OldBoard: [BoardSize][BoardSize]int{
				{2, 2, 2, 2},
				{4, 2, 2, 0},
				{2, 0, 2, 0},
				{2, 0, 0, 2},
			},
			Direction: "down",
		},
		expectedResult: moveT{
			NewBoard: [BoardSize][BoardSize]int{
				{0, 0, 0, 0},
				{2, 0, 0, 0},
				{4, 0, 2, 0},
				{4, 4, 4, 4},
			},
		},
	}
	moveTestLeft := moveTest{
		task: moveT{
			OldBoard: [BoardSize][BoardSize]int{
				{2, 2, 2, 2},
				{4, 2, 2, 0},
				{2, 0, 2, 0},
				{2, 0, 0, 2},
			},
			Direction: "left",
		},
		expectedResult: moveT{
			NewBoard: [BoardSize][BoardSize]int{
				{4, 4, 0, 0},
				{4, 4, 0, 0},
				{4, 0, 0, 0},
				{4, 0, 0, 0},
			},
		},
	}
	moveTestUp := moveTest{
		task: moveT{
			OldBoard: [BoardSize][BoardSize]int{
				{2, 2, 2, 2},
				{4, 2, 2, 0},
				{2, 0, 2, 0},
				{2, 0, 0, 2},
			},
			Direction: "up",
		},
		expectedResult: moveT{
			NewBoard: [BoardSize][BoardSize]int{
				{2, 4, 4, 4},
				{4, 0, 2, 0},
				{4, 0, 0, 0},
				{0, 0, 0, 0},
			},
		},
	}
	tests = append(tests, moveTestRight, moveTestDown, moveTestLeft, moveTestUp)
	for i := 0; i < len(tests); i = i + 1 {
		finishOffMoveTest(&tests[i])
	}
	return tests
}

func firstPassAllFieldsTest() moveTest {
	task := moveT{
		OldBoard: [BoardSize][BoardSize]int{
			{16, 8, 4, 2},
			{4, 2, 2, 0},
			{2, 4, 0, 2},
			{2, 0, 0, 0},
		},
		Direction: "right",
	}
	expectedResult := moveT{
		OldBoard: [BoardSize][BoardSize]int{
			{16, 8, 4, 2},
			{4, 2, 2, 0},
			{2, 4, 0, 2},
			{2, 0, 0, 0},
		},
		Direction: "right",
		NewBoard: [BoardSize][BoardSize]int{
			{16, 8, 4, 2},
			{0, 0, 4, 4},
			{0, 2, 4, 2},
			{0, 0, 0, 2},
		},
		NonMergeMoves: []nonMergeMoveT{
			nonMergeMoveT{{1, 0}, {1, 2}},
			nonMergeMoveT{{2, 1}, {2, 2}},
			nonMergeMoveT{{2, 0}, {2, 1}},
			nonMergeMoveT{{3, 0}, {3, 3}},
		},
		/*NewTileCandidates: []positionT{{1, 0}, {2, 0}, {3, 0},
		{3, 1}, {3, 2}, {3, 3}},*/
		MergeMoves: []mergeMoveT{
			mergeMoveT{nonMergeMoveT{{1, 2}, {1, 1}}, positionT{1, 3}, 4},
		},
		NonMovedTiles: []positionT{{0, 3}, {0, 2}, {0, 1}, {0, 0}, {2, 3}},
	}
	result := task
	result.firstPass()
	return moveTest{task, result, expectedResult}
}

func setUpMoveRightSeed() string {
	return "e9ccc20fdb924ed423ad1b46c6df43516685f4c2bc36e202ad467af1b1d1febf"
}

// This test is a simple assertion to watch the Move struct, which a lot of
// infrastructure (database, client) rely on. If the struct changes, or more
// precisely, if struct jsonification changes, then this test should fail.
//TODO rewrite Marshalling to produce a unique json representation.
/*func testJsonMarshalling(t *testing.T) {
	println("TODO implement")
	return
	move := setUpRightMoveHarness()
	marshalled, err := json.Marshal(move)
	if err != nil {
		log.Fatal(err)
	}
	expected := `{"Direction":"right","RoundNo":25,"Seed":"e9ccc20fdb924ed423ad1b46c6df43516685f4c2bc36e202ad467af1b1d1febf","OldBoard":[[16,8,4,2],[4,2,2,0],[2,4,0,2],[0,0,0,0]],"NewBoard":[[16,8,4,2],[0,0,4,4],[0,2,4,2],[0,0,0,0]],"NonMergeMoves":[[[1,0],[1,2]],[[2,0],[2,1]],[[2,1],[2,2]]],"MergeMoves":[{"Move":[[1,1],[1,2]],"Value":4}],"NonMovedTiles":[[0,0],[0,1],[0,2],[0,3]],"NewTileCandidates":[[1,0],[2,0],[3,0],[3,1],[3,2],[3,3]],"RandomTiles":[{"Position":[1,1],"Value":2}],"IsGameOver":false}`
	if string(marshalled) != expected {
		t.Error(string(marshalled) + "\n" + expected)
	}
}*/

func TestMoveTests(t *testing.T) {
	var moveTests []moveTest = make([]moveTest, 0, 10)
	moveTests = append(moveTests, firstPassAllFieldsTest())
	moveTests = append(moveTests, firstPassAllDirectionsNewBoardOnly()...)
	for _, test := range moveTests {
		result := test.check()
		if result == false {
			t.Fail()
		}
	}
}

func TestWrongDirection(t *testing.T) {
	move := moveT{
		OldBoard: [BoardSize][BoardSize]int{
			{16, 8, 4, 2},
			{4, 2, 2, 0},
			{2, 4, 0, 2},
			{2, 0, 0, 0},
		},
		Direction: "theresnosuchdirection",
	}
	move.firstPass()
	fmt.Println(move)
}
