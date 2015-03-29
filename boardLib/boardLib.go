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
// used on the front end to do board animation,
// and on the backend for future replays.
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

func (move *Move) directions(direction string) (
	initialIndex,
	smallStep,
	bigStep [2]int) {

	if direction == "left" {
		initialIndex = [2]int{0, 0}
		smallStep = [2]int{0, 1}
		bigStep = [2]int{1, -BoardSize + 1}

	} else if direction == "right" {
		initialIndex = [2]int{0, BoardSize - 1}
		smallStep = [2]int{0, -1}
		bigStep = [2]int{1, BoardSize - 1}

	} else if direction == "down" {
		initialIndex = [2]int{BoardSize - 1, 0}
		smallStep = [2]int{-1, 0}
		bigStep = [2]int{BoardSize - 1, 1}

	} else if direction == "up" {
		initialIndex = [2]int{0, 0}
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

func (move *Move) ComputeDistance() (distance BoardT, err error) {
	if !allowedMoves[move.Direction] {
		err = errors.New("has to be left, right, up or down.")
		return
	}
	initialIndex, smallStep, bigStep := move.directions(move.Direction)
	currentIndex := initialIndex
	for i := 0; i < BoardSize - 1; i++ {
		for ii := 0; ii < BoardSize - 1; ii++ {
			currentValue := move.OldBoard[currentIndex[0]][currentIndex[1]]
			currentDistance := distance[currentIndex[0]][currentIndex[1]]

			currentIndex[0] += smallStep[0]
			currentIndex[1] += smallStep[1]

			nextValue := move.OldBoard[currentIndex[0]][currentIndex[1]]
			if currentValue == 0 || currentValue == nextValue {
				distance[currentIndex[0]][currentIndex[1]] = 1 + currentDistance
			} else {
				distance[currentIndex[0]][currentIndex[1]] = currentDistance
			}
		}
		currentIndex[0] += bigStep[0]
		currentIndex[1] += bigStep[1]
	}
	return
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
