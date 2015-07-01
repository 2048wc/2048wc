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

// import "encoding/json"
import "fmt"

// The size of the board. Most likely will always stay 4,
// as this is what seems to be most playable.
const BoardSize = 4

var allowedMoves map[string]bool = map[string]bool{"left": true, "right": true,
	"up": true, "down": true}

type BoardT [BoardSize][BoardSize]int
type positionT [2]int
type positionValueT struct {
	Position positionT
	Value    int
}
type moveT [2]positionT
type mergeMoveT struct {
	From  moveT
	To    positionT
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
	//You are responsible for initialising those
	Direction string
	RoundNo   int
	Seed      string
	OldBoard  BoardT
	//These are going to be computed by the move pipeline
	NewBoard          BoardT
	NonMergeMoves     []moveT
	MergeMoves        []mergeMoveT
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
func CreateMove(oldBoard BoardT,
	direction string,
	roundNo int,
	seed string) (move Move) {
	move.Direction = direction
	move.Seed = seed
	move.OldBoard = oldBoard
	move.RoundNo = roundNo
	move.NonMergeMoves = make([]moveT, 0, BoardSize*BoardSize)
	move.MergeMoves = make([]mergeMoveT, 0, BoardSize*BoardSize)
	move.NonMovedTiles = make([]positionT, 0, BoardSize*BoardSize)
	move.NewTileCandidates = make([]positionT, 0, BoardSize*BoardSize)
	move.RandomTiles = make([]positionValueT, 0, BoardSize*BoardSize)
	return
}

func (board *BoardT) get(position positionT) int {
	return board[position[0]][position[1]]
}

func (board *BoardT) set(position positionT, value int) {
	board[position[0]][position[1]] = value
	return
}

func (board *BoardT) add(position positionT, value int) int {
	board[position[0]][position[1]] += value
	return board[position[0]][position[1]]
}

//TODO implement to exlude wrong seeds, wrong directions, etc.
func (move *Move) validateBoardInitialisation() {
	return
}

//TODO checks if game is over
func (move *Move) isGameOver() bool {
	return false
}

func invertDirection(direction string) string {
	if direction == "up" {
		return "down"
	}
	if direction == "down" {
		return "up"
	}
	if direction == "left" {
		return "right"
	}
	if direction == "right" {
		return "left"
	}
	return ""
}

type IterationStateMachine struct {
	// These are assigned by setDirections
	smallStepForward  [2]int
	smallStepBackward [2]int
	bigStep           [2]int
	// These are modified in every step of iteration
	distance                int
	mergeHopefulIndex       positionT
	isMergeHopeful          bool
	isLast                  bool
	isHopefulUnmoved        bool
	mergeHopefulDestination positionT
	// This is set by setDirections and then modified every step of iteration
	currentIndex positionT
}

func (ism *IterationStateMachine) setDirections(direction string) {
	if direction == "right" {
		ism.currentIndex = positionT{0, BoardSize - 1}
		ism.smallStepForward = [2]int{0, 1}
		ism.smallStepBackward = [2]int{0, -1}
		ism.bigStep = [2]int{1, BoardSize - 1}
	} else if direction == "left" {
		ism.currentIndex = positionT{0, 0}
		ism.smallStepForward = [2]int{0, -1}
		ism.smallStepBackward = [2]int{0, 1}
		ism.bigStep = [2]int{1, -BoardSize + 1}
	} else if direction == "up" {
		ism.currentIndex = positionT{0, 0}
		ism.smallStepForward = [2]int{-1, 0}
		ism.smallStepBackward = [2]int{1, 0}
		ism.bigStep = [2]int{-BoardSize + 1, 1}
	} else if direction == "down" {
		ism.currentIndex = positionT{BoardSize - 1, 0}
		ism.smallStepForward = [2]int{1, 0}
		ism.smallStepBackward = [2]int{-1, 0}
		ism.bigStep = [2]int{BoardSize - 1, 1}
	}
	return
}

func (position *positionT) nextPosition(stepSize [2]int, numberOfSteps int) {
	position[0] = position[0] + stepSize[0]*numberOfSteps
	position[1] = position[1] + stepSize[1]*numberOfSteps
}

func (newboard *BoardT) incrementFromBySteps(oldboard *BoardT,
	incrementFrom positionT, step [2]int, numberOfSteps int) (int, positionT) {
	value := oldboard[incrementFrom[0]][incrementFrom[1]]
	row := incrementFrom[0] + step[0]*numberOfSteps
	column := incrementFrom[1] + step[1]*numberOfSteps
	newboard[row][column] += value
	var tusia = 0
	tusia = tusia
	return newboard[row][column], positionT{row, column}
}

func (move *Move) resolveBoardAtIndex(ism *IterationStateMachine) {
	giveHope := func() {
		_, position := move.NewBoard.incrementFromBySteps(&move.OldBoard,
			ism.currentIndex, ism.smallStepForward, ism.distance)
		ism.mergeHopefulIndex = ism.currentIndex
		ism.mergeHopefulDestination = position
		ism.isMergeHopeful = true
		if ism.distance == 0 {
			ism.isHopefulUnmoved = true
		} else {
			ism.isHopefulUnmoved = false
		}
	}
	abandonHope := func() {
		ism.isMergeHopeful = false
		ism.isHopefulUnmoved = false
		ism.mergeHopefulDestination = positionT{-1, -1}
		ism.mergeHopefulIndex = positionT{-1, -1}
	}
	dispatchLoser := func() {
		if ism.isHopefulUnmoved == true {
			move.NonMovedTiles = append(move.NonMovedTiles, ism.mergeHopefulIndex)
		} else {
			move.NonMergeMoves = append(move.NonMergeMoves,
				moveT{ism.mergeHopefulIndex, ism.mergeHopefulDestination})
		}
	}
	completeMerge := func() {
		ism.distance += 1
		value, position := move.NewBoard.incrementFromBySteps(&move.OldBoard,
			ism.currentIndex, ism.smallStepForward, ism.distance)
		fromMoves := moveT{ism.mergeHopefulIndex, ism.currentIndex}
		mergeMove := mergeMoveT{From: fromMoves,
			To:    position,
			Value: value}
		move.MergeMoves = append(move.MergeMoves, mergeMove)
	}
	currentValue := move.OldBoard.get(ism.currentIndex)
	if currentValue == 0 {
		ism.distance += 1
	} else {
		if ism.isMergeHopeful {
			if move.OldBoard.get(ism.mergeHopefulIndex) == currentValue {
				completeMerge()
				abandonHope()
			} else {
				dispatchLoser()
				giveHope()
			}
		} else {
			giveHope()
		}
	}
	if ism.isLast {
		if ism.isMergeHopeful {
			dispatchLoser()
			abandonHope()
		}
		ism.isLast = false
	}
}

//TODO implement to satisfy the TDD test case
func (move *Move) ExecuteMove() error {
	var ism IterationStateMachine
	ism.setDirections(move.Direction)
	for i := 0; i < BoardSize; i++ {
		ism.distance = 0
		move.resolveBoardAtIndex(&ism)
		for ii := 0; ii < BoardSize-1; ii++ {
			ism.currentIndex.nextPosition(ism.smallStepBackward, 1)
			if ii == BoardSize-2 {
				ism.isLast = true
			}
			move.resolveBoardAtIndex(&ism)
		}
		/*if i != BoardSize {*/
			ism.currentIndex.nextPosition(ism.bigStep, 1)
		/*}*/
	}
	move.RoundNo += 1
	return nil
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
