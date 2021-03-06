package boardLib

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

import (
	"API2048"
	"fmt"
	"os"
	"reflect"
	"testing"
)

var moveCreator API2048.MoveCreator = MoveCreator{}

// test that the board is an n by n iterable where all elements are 0
func TestBoardInitialised(t *testing.T) {
	var board boardT
	for i := 0; i < API2048.BoardSize; i++ {
		for j := 0; j < API2048.BoardSize; j++ {
			if board[i][j] != 0 {
				t.Fail()
			}
		}
	}
}

// we will mostly likely only use board size of 4. If one day
// we decide otherwise, we can update the test
func TestBoardSize(t *testing.T) {
	if API2048.BoardSize != 4 {
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
			OldBoard: [API2048.BoardSize][API2048.BoardSize]int{
				{2, 2, 2, 2},
				{4, 2, 2, 0},
				{2, 0, 2, 0},
				{2, 0, 0, 2},
			},
			Direction: "right",
		},
		expectedResult: moveT{
			NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
				{0, 0, 4, 4},
				{0, 0, 4, 4},
				{0, 0, 0, 4},
				{0, 0, 0, 4},
			},
		},
	}
	moveTestDown := moveTest{
		task: moveT{
			OldBoard: [API2048.BoardSize][API2048.BoardSize]int{
				{2, 2, 2, 2},
				{4, 2, 2, 0},
				{2, 0, 2, 0},
				{2, 0, 0, 2},
			},
			Direction: "down",
		},
		expectedResult: moveT{
			NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
				{0, 0, 0, 0},
				{2, 0, 0, 0},
				{4, 0, 2, 0},
				{4, 4, 4, 4},
			},
		},
	}
	moveTestLeft := moveTest{
		task: moveT{
			OldBoard: [API2048.BoardSize][API2048.BoardSize]int{
				{2, 2, 2, 2},
				{4, 2, 2, 0},
				{2, 0, 2, 0},
				{2, 0, 0, 2},
			},
			Direction: "left",
		},
		expectedResult: moveT{
			NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
				{4, 4, 0, 0},
				{4, 4, 0, 0},
				{4, 0, 0, 0},
				{4, 0, 0, 0},
			},
		},
	}
	moveTestUp := moveTest{
		task: moveT{
			OldBoard: [API2048.BoardSize][API2048.BoardSize]int{
				{2, 2, 2, 2},
				{4, 2, 2, 0},
				{2, 0, 2, 0},
				{2, 0, 0, 2},
			},
			Direction: "up",
		},
		expectedResult: moveT{
			NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
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
		OldBoard: [API2048.BoardSize][API2048.BoardSize]int{
			{16, 8, 4, 2},
			{4, 2, 2, 0},
			{2, 4, 0, 2},
			{2, 0, 0, 0},
		},
		Direction: "right",
	}
	expectedResult := moveT{
		OldBoard: [API2048.BoardSize][API2048.BoardSize]int{
			{16, 8, 4, 2},
			{4, 2, 2, 0},
			{2, 4, 0, 2},
			{2, 0, 0, 0},
		},
		Direction: "right",
		NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
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
		OldBoard: [API2048.BoardSize][API2048.BoardSize]int{
			{16, 8, 4, 2},
			{4, 2, 2, 0},
			{2, 4, 0, 2},
			{2, 0, 0, 0},
		},
		Direction: "theresnosuchdirection",
	}
	move.firstPass()
}

type magic struct {
	Magic1 int
	Magic2 string
	Magic3 [2]int
}

func TestMarshalWithoutFields(t *testing.T) {
	var magic = magic{0, "tusia", [2]int{0, 1}}
	if json, status := marshalExcludeFields(magic,
		map[string]bool{"Magic1": true}); status != nil {
		t.Error(`encoding/json failed`)
	} else {
		if string(json) != `{"Magic2":"tusia","Magic3":[0,1]}` {
			t.Error("is: ", json,
				"should be: ", `{"Magic2":"tusia","Magic3":[0,1]}`)
		}
	}
}

// 12259964326927110866866776217202473468949912977468817261 is
// 7fffffffffffffffffffffffffffffffffffffffffff6d in hex.
func TestRandom(t *testing.T) {
	move := moveCreator.CreateMove().(*moveT)
	move.InitMove([API2048.BoardSize][API2048.BoardSize]int{
		{16, 8, 4, 2},
		{4, 2, 2, 0},
		{2, 4, 0, 2},
		{2, 0, 0, 0},
	}, "left", 21, "7fffffffffffffffffffffffffffffffffffffffffff6d")
	seedPlusOne := moveCreator.CreateMove().(*moveT)
	seedPlusOne.InitMove([API2048.BoardSize][API2048.BoardSize]int{
		{16, 8, 4, 2},
		{4, 2, 2, 0},
		{2, 4, 0, 2},
		{2, 0, 0, 0},
	}, "left", 20, "7fffffffffffffffffffffffffffffffffffffffffff6e")
	moveString := move.Seed.String()
	if moveString != "12259964326927110866866776217202473468949912977468817261" {
		t.Error("wrong hex->int bigInt conversion.", move.Seed.String())
	}
	if move.randInt(100, false) != seedPlusOne.randInt(100, false) {
		t.Error("roundNo doesn't have equal input with the seed")
	}
	if move.randInt(65536, false) != move.randInt(65536, false) {
		t.Error("should be indempotent")
	}
	seedPlusOne.RoundNo = 21
	if move.randInt(19911993, false) != seedPlusOne.randInt(19911993, true) {
		t.Error("can't do previous.")
	}
}

func TestSecondPass(t *testing.T) {
	noGameOverColumns := moveT{
		NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10, 11, 12},
			{13, 14, 15, 12},
		},
	}
	noGameOverColumns.secondPass()
	if noGameOverColumns.IsGameOver {
		t.Error("why game over")
	}
	if len(noGameOverColumns.NewTileCandidates) != 0 {
		t.Error("should be empty")
	}

	noGameOverRows := moveT{
		NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10, 11, 12},
			{13, 14, 15, 15},
		},
	}
	noGameOverRows.secondPass()
	if noGameOverRows.IsGameOver {
		t.Error("why game over")
	}
	if len(noGameOverRows.NewTileCandidates) != 0 {
		t.Error("should be empty")
	}

	noGameOverEmpties := moveT{
		NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
			{1, 2, 3, 4},
			{0, 6, 0, 8},  // positionT{1, 0}, positionT{1, 2}
			{9, 0, 11, 0}, // positionT{2, 1}, positionT{2, 3}
			{13, 14, 15, 15},
		},
	}

	noGameOverEmpties.secondPass()
	if noGameOverEmpties.IsGameOver {
		t.Error("why game over")
	}
	if reflect.DeepEqual(noGameOverEmpties.NewTileCandidates,
		[]positionT{{1, 0}, {1, 2}, {2, 1}, {2, 3}},
	) == false {
		t.Error("should be empty")
	}
	gameOver := moveT{
		NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
			{1, 2, 3, 4},
			{4, 6, 1, 8},
			{9, 1, 11, 1},
			{13, 14, 15, 8},
		},
	}
	gameOver.secondPass()
	if !gameOver.IsGameOver {
		t.Error("should game over")
	}
	gameOver2 := moveT{
		NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
			{1, 2, 1, 2},
			{2, 1, 2, 1},
			{1, 2, 1, 2},
			{2, 1, 2, 1},
		},
	}
	gameOver2.secondPass()
	if !gameOver2.IsGameOver {
		t.Error("should game over")
	}
	topLeft := moveT{
		NewBoard: [API2048.BoardSize][API2048.BoardSize]int{
			{3, 3, 1, 2},
			{2, 1, 2, 1},
			{1, 2, 1, 2},
			{2, 1, 2, 1},
		},
	}
	topLeft.secondPass()
	if topLeft.IsGameOver {
		t.Error("why game over")
	}
}

func TestInternalView(t *testing.T) {
	moveI := moveCreator.CreateMove()
	moveI.InitMove([API2048.BoardSize][API2048.BoardSize]int{}, "", 0, "7fffffffffffffffffffffffffffffffffffffffffff6e")
	expected := `{"Direction":"","RoundNo":0,"Seed":"7fffffffffffffffffffffffffffffffffffffffffff6e","OldBoard":[[0,0,0,0],[0,0,0,0],[0,0,0,0],[0,0,0,0]],"NewBoard":[[0,0,0,0],[0,0,0,0],[0,0,0,0],[0,0,0,0]],"NonMergeMoves":[],"MergeMoves":[],"NonMovedTiles":[],"NewTileCandidates":[],"IsGameOver":false,"RandomTiles":[]}`
	actual := moveI.InternalView()
	if actual != expected {
		t.Error("\nactual: --\n expected: --\n", actual, "\n", expected)
	}
}

func TestHashing(t *testing.T) {
	movea := moveCreator.CreateMove().(*moveT)
	// hashing a string "a", also known as a byte 0x61 or simply 0110 0001
	// SHA256 of this string is ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb
	// In decimal it is 91634880152443617534842621287039938041581081254914058002978601050179556493499
	// the remainder from division of this number by 1 million is 493499
	// could be a lucky coincidence ;)
	movea.InitMove([API2048.BoardSize][API2048.BoardSize]int{}, "", 0, "61")
	if movea.randInt(1000*1000, false) != 493499 {
		t.Error(movea.randInt(1000*1000, false), "!=", 493499)
	}
}

func TestUnmarshalSeed(t *testing.T) {
	movea := moveCreator.CreateMove().(*moveT)
	jsona := `{"Direction":"","RoundNo":0,"Seed":"7fffffffffffffffffffffffffffffffffffffffffff6e","OldBoard":[[1,0,0,0],[0,0,0,0],[0,0,0,0],[0,0,0,0]],"NewBoard":[[0,0,0,0],[0,0,0,0],[0,0,0,0],[0,0,0,0]],"NonMergeMoves":[],"MergeMoves":[],"NonMovedTiles":[],"NewTileCandidates":[],"IsGameOver":false,"RandomTiles":[]}`
	movea.ParseMove(jsona)
	if movea.InternalView() != jsona {
		fmt.Println("")
		fmt.Println(movea.InternalView())
		fmt.Println(jsona)
		t.Error("mismatch")
	}
}

func TestResolveMove(t *testing.T) {
	var move API2048.Move
	move = moveCreator.CreateMove()
	move.InitMove([API2048.BoardSize][API2048.BoardSize]int{
		{0, 0, 0, 2},
		{0, 2, 0, 0},
		{0, 0, 4, 8},
		{0, 8, 64, 8},
	}, "up", 24, "7fffffffffffffffffffffffffffffffffffffffffff6e")
	move.ResolveMove()
	expected := `{"Direction":"up","RoundNo":25,"Seed":"7fffffffffffffffffffffffffffffffffffffffffff6e","OldBoard":[[0,0,0,2],[0,2,0,0],[0,0,4,8],[0,8,64,8]],"NewBoard":[[0,2,4,2],[0,8,64,16],[0,0,0,0],[0,0,0,2]],"NonMergeMoves":[[[1,1],[0,1]],[[3,1],[1,1]],[[2,2],[0,2]],[[3,2],[1,2]]],"MergeMoves":[{"From":[[2,3],[3,3]],"To":[1,3],"Value":16}],"NonMovedTiles":[[0,3]],"NewTileCandidates":[[0,0],[1,0],[2,0],[2,1],[2,2],[2,3],[3,0],[3,1],[3,2],[3,3]],"IsGameOver":false,"RandomTiles":[{"Position":[3,3],"Value":2}]}`
	if move.InternalView() != expected {
		fmt.Println(move.InternalView())
		fmt.Println(expected)
		t.Error("move is not properly resolved")
	}
}

func TestFullPipeline(t *testing.T) {
	movea := moveCreator.CreateMove().(*moveT)
	jsona := `{"Direction":"up","RoundNo":24,"Seed":"7fffffffffffffffffffffffffffffffffffffffffff6e","OldBoard":[[0,0,0,2],[0,2,0,0],[0,0,4,8],[0,8,64,8]],"NewBoard":[[0,2,4,2],[0,8,64,16],[0,0,0,0],[2,0,0,0]],"NonMergeMoves":[[[1,1],[0,1]],[[3,1],[1,1]],[[2,2],[0,2]],[[3,2],[1,2]]],"MergeMoves":[{"From":[[2,3],[3,3]],"To":[1,3],"Value":16}],"NonMovedTiles":[[0,3]],"NewTileCandidates":[[0,0],[1,0],[2,0],[2,1],[2,2],[2,3],[3,0],[3,1],[3,2],[3,3]],"IsGameOver":false,"RandomTiles":[{"Position":[3,0],"Value":2}]}`
	movea.ParseMove(jsona)
	moveb := movea.CreateNextMove()
	if moveb.GetRoundNo() != 24 {
		t.Error("do not round up")
	}
	if moveb.GetSeed() != "7fffffffffffffffffffffffffffffffffffffffffff6e" {
		t.Error("seed stays the same")
	}
	moveb.SetDirection("up")
	moveb.ResolveMove()
	if moveb.GetRoundNo() != 25 {
		t.Error("do round up")
	}
	expectedInternalView := `{"Direction":"up","RoundNo":25,"Seed":"7fffffffffffffffffffffffffffffffffffffffffff6e","OldBoard":[[0,2,4,2],[0,8,64,16],[0,0,0,0],[2,0,0,0]],"NewBoard":[[2,2,4,2],[0,8,64,16],[0,0,0,0],[0,0,2,0]],"NonMergeMoves":[[[3,0],[0,0]]],"MergeMoves":[],"NonMovedTiles":[[0,1],[1,1],[0,2],[1,2],[0,3],[1,3]],"NewTileCandidates":[[1,0],[2,0],[2,1],[2,2],[2,3],[3,0],[3,1],[3,2],[3,3]],"IsGameOver":false,"RandomTiles":[{"Position":[3,2],"Value":2}]}`
	if moveb.InternalView() != expectedInternalView {
		fmt.Println(moveb.InternalView())
		fmt.Println(expectedInternalView)
		t.Error("error in moveb")
	}
}

func TestInit(t *testing.T) {
	var move API2048.Move
	move = moveCreator.CreateMove()
	move.InitFirstMove()
	newTiles := 0
	for i := 0; i < API2048.BoardSize; i++ {
		for j := 0; j < API2048.BoardSize; j++ {
			if move.(*moveT).OldBoard[i][j] == NewTileValue {
				newTiles += 1
			}
		}
	}
	if newTiles != 2 {
		t.Error("need 2 new tiles to start with")
	}
}

func TestGameOver(t *testing.T) {
	move := moveCreator.CreateMove()
	board := [API2048.BoardSize][API2048.BoardSize]int{
		{4, 8, 64, 2},
		{8, 4, 8, 4},
		{4, 32, 4, 2},
		{4, 128, 16, 0},
	}
	direction := "right"
	move.InitMove(board, direction, 10, "a")
	move.ResolveMove()
	if !move.GetGameOver() {
		t.Error("this needs to game over.")
	}
	if move.GetRoundNo() != 11 {
		t.Error("bad round number.")
	}

}
