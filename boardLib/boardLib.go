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

import "fmt"
import "errors"

// The size of the board. Most likely will always stay 4,
// as this is what seems to be most playable.
const BoardSize = 4

type BoardT [BoardSize][BoardSize]int
type positionT [2]int
type positionValueT struct {
	Position positionT
	Value    int
}
type moveT [2]positionT
type moveValueT struct {
	Move  moveT
	Value int
}

// This is the struct, which is going to be both
// used on the frontend to do board animation,
// and on the backend for future replays.
// both fronted and backend json views of this struct
// are going to look the same, apart from
// dropping the Seed field from the client view
// for security reasons.
type Move struct {
	//These are created by the constructor
	Direction string
	RoundNo   int
	Seed      string
	OldBoard  BoardT
	//These are computed by the move pipeline
	NewBoard          BoardT
	NonMergeMoves     []moveT
	MergeMoves        []moveValueT
	NonMovedTiles     []positionT
	NewTileCandidates []positionT
	RandomTiles       []positionValueT
	IsGameOver        bool
}

// Direction is either "left", "right", "up" or "down"
// RoundNo is also a score
// Seed represents a 256 bit number encoded as 64 character long hex number.
// OldBoard represents a board to evaluate
// TODO validate values
func CreateMove(oldBoard BoardT, direction string, roundNo int, seed string) (move Move) {
	move.Direction = direction
	move.Seed = seed
	move.OldBoard = oldBoard
	return
}

//TODO implement to satisfy the TDD test case
func (move *Move) ExecuteMove() {
	return
}

func directions(direction string) (
	initialIndex positionT,
	smallStep [2]int,
	bigStep [2]int) {

	if direction == "left" {
		initialIndex = positionT{0, 0}
		smallStep = [2]int{0, 1}
		bigStep = [2]int{1, -BoardSize + 1}

	} else if direction == "right" {
		initialIndex = positionT{0, BoardSize - 1}
		smallStep = [2]int{0, -1}
		bigStep = [2]int{1, BoardSize - 1}

	} else if direction == "down" {
		initialIndex = positionT{BoardSize - 1, 0}
		smallStep = [2]int{-1, 0}
		bigStep = [2]int{BoardSize - 1, 1}

	} else if direction == "up" {
		initialIndex = positionT{0, 0}
		smallStep = [2]int{1, 0}
		bigStep = [2]int{-BoardSize + 1, 1}

	} else {
		panic("has to be left, right, up or down.")
	}
	return
}

var allowedMoves map[string]bool

func init() {
	allowedMoves = make(map[string]bool)
	allowedMoves["left"] = true
	allowedMoves["right"] = true
	allowedMoves["up"] = true
	allowedMoves["down"] = true
}

func boardValue(board BoardT, position positionT) int {
	return board[position[0]][position[1]]
}

func (position *positionT) nextPosition(stepSize [2]int, numberOfSteps int) {
	position[0] = position[0] + stepSize[0]*numberOfSteps
	position[1] = position[1] + stepSize[1]*numberOfSteps
}

//a struct to represent a unit of iteration over the board.
type boardIndexT struct {
	currentIndex positionT
	nextIndex    positionT
	currentValue int
	nextValue    int
	newPass      bool
}

func boardIterator(board BoardT, direction string) (
	iterator func() (boardIndex boardIndexT, err error)) {

	var boardIndex boardIndexT
	initialIndex, smallStep, bigStep := directions(direction)
	boardIndex.nextIndex = initialIndex

	var iter int
	var innerIter int

	var makeSteps = func() (err error) {
		if innerIter < BoardSize-1 {
			innerIter += 1
			boardIndex.nextIndex.nextPosition(smallStep, 1)
			boardIndex.newPass = false
		} else if iter < BoardSize-1 {
			innerIter = 0
			iter += 1
			boardIndex.nextIndex.nextPosition(bigStep, 1)
			boardIndex.newPass = true
		} else {
			err = errors.New("iterator out of bonds")
		}
		return
	}
	iterator = func() (boardIndexT, error) {
		boardIndex.currentIndex = boardIndex.nextIndex
		boardIndex.currentValue = boardValue(board, boardIndex.currentIndex)
		err := makeSteps()
		boardIndex.nextValue = boardValue(board, boardIndex.nextIndex)
		return boardIndex, err
	}
	return
}

func (move *Move) computeDistance() (BoardT, error) {
	var distance BoardT
	iterator := boardIterator(move.OldBoard, move.Direction)
	for i := 0; i < BoardSize*BoardSize-1; i++ {
		ii, err := iterator()
		if err != nil {
			return distance, err
		}
		currentDistance := distance[ii.currentIndex[0]][ii.currentIndex[1]]
		if ii.newPass == true {
			currentDistance = 0
		}
		if ii.currentValue == 0 || ii.currentValue == ii.nextValue {
			distance[ii.nextIndex[0]][ii.nextIndex[1]] = 1 + currentDistance
		} else {
			distance[ii.nextIndex[0]][ii.nextIndex[1]] = currentDistance
		}
	}
	return distance, nil
}

// TODO implement, and cover with a test case,
// perhaps smarter than the current marshalling test case
func (move *Move) Marshal() (jsonified string) {
	return ""
}

// This might be useful for logging/ command line interface.
// For the database we'll have to use a different function that's
// going to return json.
func PrintBoard(board BoardT) {
	for _, row := range board {
		for _, value := range row {
			fmt.Print(value, " ")
		}
		fmt.Println("")
	}
}
