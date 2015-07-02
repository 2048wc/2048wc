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
//import "fmt"

// The size of the board. Most likely will always stay 4,
// as this is what seems to be most playable.
const BoardSize = 4

var allowedMoves map[string]bool = map[string]bool{"left": true, "right": true,
	"up": true, "down": true}

type boardT [BoardSize][BoardSize]int
type positionT [2]int
type positionValueT struct {
	Position positionT
	Value    int
}

type nonMergeMoveT [2]positionT
type mergeMoveT struct {
	From  nonMergeMoveT
	To    positionT
	Value int
}

func CreateMove() Move {
	return &moveT{}
}

type Move interface {
	InitFirstMove()
	InitNextMove(oldBoard [BoardSize][BoardSize]int, direction string, roundNo int, seed string)
	ValidateMove() (success bool)
	ResolveMove() 
	InternalView() (json string)
	ExternalView() (json string)
	ParseMove(string) 
	GetRoundNo() int
}

// The json view of this struct is going to be both
// used on the frontend to do board animation,
// and on the backend for future replays.
// both fronted and backend json views of this struct
// are going to look the same, apart from
// dropping the Seed field from the client view
// for security reasons.
type moveT struct {
	//You are responsible for initialising those
	Direction string
	RoundNo   int
	Seed      string
	OldBoard  boardT
	// Resolved by firstPass
	NewBoard          boardT
	NonMergeMoves     []nonMergeMoveT
	MergeMoves        []mergeMoveT
	NonMovedTiles     []positionT
	// Resolved by secondPass
	NewTileCandidates []positionT
	IsGameOver        bool
	// Resolved by addRandomTile
	RandomTiles       []positionValueT
}


func (move *nonMergeMoveT) ExternalView() string{
	return ""
}

func (move *moveT) GetRoundNo() int{
	return move.RoundNo
}

func (move *moveT) ParseMove(string){
	return 
}

func (move *moveT) ValidateMove() bool{
	return true
}

// Direction is either "left", "right", "up" or "down"
// RoundNo is also a score
// Seed represents a 256 bit number encoded as 64 character long hex number.
// OldBoard represents a board to evaluate
// TODO validate values
func (move *moveT) InitNextMove(oldBoard [BoardSize][BoardSize]int,
	direction string,
	roundNo int,
	seed string) {
	move.Direction = direction
	move.Seed = seed
	move.OldBoard = oldBoard
	move.RoundNo = roundNo
	move.initMoveCollections()
	return
}

func (move *moveT) InitFirstMove() {
	move.initMoveCollections()
	return
}

func (move *moveT) InternalView() string{
	return ""
}

func (move *moveT) ExternalView() string{
	return ""
}

func (move *moveT) initMoveCollections(){	
	move.NonMergeMoves = make([]nonMergeMoveT, 0, BoardSize*BoardSize)
	move.MergeMoves = make([]mergeMoveT, 0, BoardSize*BoardSize)
	move.NonMovedTiles = make([]positionT, 0, BoardSize*BoardSize)
	move.NewTileCandidates = make([]positionT, 0, BoardSize*BoardSize)
	move.RandomTiles = make([]positionValueT, 0, BoardSize*BoardSize)
}	

func (board *boardT) get(position positionT) int {
	return board[position[0]][position[1]]
}

func (board *boardT) set(position positionT, value int) {
	board[position[0]][position[1]] = value
	return
}


//TODO implement to exlude wrong seeds, wrong directions, etc.
func (move *moveT) validateBoardInitialisation() {
	return
}

//TODO checks if game is over
func (move *moveT) isGameOver() bool {
	return false
}

type iterationStateMachine struct {
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

func (ism *iterationStateMachine) setDirections(direction string) {
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

func (newboard *boardT) incrementFromBySteps(oldboard *boardT,
	incrementFrom positionT, step [2]int, numberOfSteps int) (int, positionT) {
	value := oldboard.get(incrementFrom)
	row := incrementFrom[0] + step[0]*numberOfSteps
	column := incrementFrom[1] + step[1]*numberOfSteps
	newboard[row][column] += value
	return newboard[row][column], positionT{row, column}
}

func (move *moveT) resolveBoardAtIndex(ism *iterationStateMachine) {
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
				nonMergeMoveT{ism.mergeHopefulIndex, ism.mergeHopefulDestination})
		}
	}
	completeMerge := func() {
		ism.distance += 1
		value, position := move.NewBoard.incrementFromBySteps(&move.OldBoard,
			ism.currentIndex, ism.smallStepForward, ism.distance)
		fromMoves := nonMergeMoveT{ism.mergeHopefulIndex, ism.currentIndex}
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

func (move *moveT) firstPass(){
	var ism iterationStateMachine
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
}

//TODO implement to satisfy the TDD test case
func (move *moveT) ResolveMove() {
	move.firstPass()
	move.RoundNo += 1
	return
}

// This might be useful for logging/ command line interface.
// For the database we'll have to use a different function that's
// going to return json.
/*func PrintBoard(board boardT) {
	for _, row := range board {
		for _, value := range row {
			fmt.Print(value, " ")
		}
		fmt.Println("")
	}
}*/