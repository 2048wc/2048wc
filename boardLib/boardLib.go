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

import "encoding/json"
import "encoding/hex"
import "reflect"
import "math/big"
import "fmt"
//import "errors"

import "crypto/sha256"

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

// Creates an empty Move. The move has to be then initialised using
func CreateMove() Move {
	return &moveT{}
}

type Move interface {
	// This initialises all the fields with sensible defaults for the first
	// move. This includes 2 random tiles, a random seed, a non-zero round
	// number, a valid direction and non-nil pointers for all internal data
	// structures. After initialisation it is safe to call other functions on
	// this struct.
	InitFirstMove()

	// Like InitFirstMove, but with more control over the initial conditions of
	InitMove(oldBoard [BoardSize][BoardSize]int,
		direction string, roundNo int, seed string)

	// Init the board using a json passed in as a string. Call ValidateMove
	// after this function.
	ParseMove(json string)

	// Sets direction of this move. Allowed directions are left, right,
	// down, up. Call ValidateMove afterwards to make sure that the struct is
	// still in acceptable state.
	SetDirection(string)

	// Returns a new move newMove.OldBoard := oldMove.NewBoard, roundNo += 1,
	// and the seed carried over.
	CreateNextMove() Move

	// Evolves board one step forward. Can put errors in ValidateMove.
	// Call ValidateMove before and afterwards.
	ResolveMove()

	// Checks if move fields satisfy some basic constraints. Should be
	// called before ResolveMove. Returns a map of field names onto errors for
	// incorrectly valued fields. Both Resolved and Unresolved moves satisfy
	// validation. Resolved moves return isResolved true. Move which doesn't
	// progress the board is detected as an error in NewBoard.
	ValidateMove() (isResolved bool, errors map[string]error)

	// Internal json representation of the struct. Do not show to the client!
	// Exports confidential fields.
	InternalView() (json string)

	// External json representation. Safe to share with the client.
	ExternalView() (json string)

	// Get Round Number, which is simultanously a score.
	GetRoundNo() int

	// Get Seed
	GetSeed() string

	// Get the status of the game.
	GetGameOver() bool
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
	Seed      big.Int
	OldBoard  boardT
	// Resolved by firstPass
	NewBoard      boardT
	NonMergeMoves []nonMergeMoveT
	MergeMoves    []mergeMoveT
	NonMovedTiles []positionT
	// Resolved by secondPass
	NewTileCandidates []positionT
	IsGameOver        bool
	RandomTiles       []positionValueT
}

type internalMoveT struct {
	//You are responsible for initialising those
	Direction string
	RoundNo   int
	Seed      string
	OldBoard  boardT
	// Resolved by firstPass
	NewBoard      boardT
	NonMergeMoves []nonMergeMoveT
	MergeMoves    []mergeMoveT
	NonMovedTiles []positionT
	// Resolved by secondPass
	NewTileCandidates []positionT
	IsGameOver        bool
	RandomTiles       []positionValueT
}

func (move *moveT) GetGameOver() bool {
	return move.IsGameOver
}

func (move *moveT) SetDirection(string) {
	return
}

func (move *moveT) ExternalView() string {
	jsona, err := marshalExcludeFields(*move, map[string]bool{"Seed": true})
	//TODO put err into a map with errors
	if err != nil{
		return ""
	}
	return string(jsona)
}

func (move *moveT) InternalView() string {
	var jsona []byte
	value := reflect.ValueOf(*move)
	typa := reflect.TypeOf(*move)
	size := value.NumField()
	jsona = append(jsona, '{')
	for i := 0; i < size; i++ {
		structValue := value.Field(i)
		fieldName := typa.Field(i).Name
		var marshalStatus error
		var marshalledField []byte
		if fieldName == "Seed"{
			marshalledField, marshalStatus = json.Marshal(
				move.GetSeed());
			if marshalStatus != nil {
				return ""
			}
		} else{
			marshalledField, marshalStatus = json.Marshal(
			(structValue).Interface())
			if marshalStatus != nil {
			return ""
			}
		}
		jsona = append(jsona, '"')
		jsona = append(jsona, []byte(fieldName)...)
		jsona = append(jsona, '"')
		jsona = append(jsona, ':')
		jsona = append(jsona, (marshalledField)...)
		if i != size - 1{
			jsona = append(jsona, ',')
		}
	}
	jsona = append(jsona, '}')
	return string(jsona)
}

func (move *moveT) GetRoundNo() int {
	return move.RoundNo
}

func (move *moveT) ParseMove(jsona string) {
	var internalMoveStruct internalMoveT
	if errora := json.Unmarshal([]byte(jsona), &internalMoveStruct); errora != nil {
		//TODO set error flags
		fmt.Println(errora)
		panic(errora)
		return
	}
	// repr stands for representation, as in-memory representation. That's to
	// differentiate it from internal representation.
	
	internalMove := reflect.ValueOf(internalMoveStruct)
	internalMoveType := reflect.TypeOf(internalMoveStruct)
	reprMove := reflect.ValueOf(move)
	//reprMoveType := reflect.TypeOf(*move)
	size := internalMove.NumField()
	if size != reprMove.Elem().NumField() {
		// TODO set error flags
		panic("misaligned structs")
	}
	for i := 0; i < size; i++ {
		internalField := internalMove.Field(i)
		internalFieldName := internalMoveType.Field(i).Name
		reprField := reprMove.Elem().FieldByName(internalFieldName)
		if internalFieldName == "Seed" {
			var bigInt big.Int
			bigInt.SetString(internalField.Interface().(string), 16)
			move.Seed = bigInt
		} else {
			reprField.Set(internalField)
		}
	}
}

func (move *moveT) ValidateMove() (bool, map[string]error) {
	return true, nil
}

// Direction is either "left", "right", "up" or "down"
// RoundNo is also a score
// Seed represents a 256 bit number encoded as 64 character long hex number.
// OldBoard represents a board to evaluate
// TODO validate values
func (move *moveT) InitMove(oldBoard [BoardSize][BoardSize]int,
	direction string, roundNo int, seed string) {
	move.Direction = direction
	move.Seed.SetString(seed, 16)
	move.OldBoard = oldBoard
	move.RoundNo = roundNo
	move.initMoveCollections()
}

func (move *moveT) CreateNextMove() Move {
	nextMove := &moveT{}
	nextMove.initMoveCollections()
	nextMove.Seed = move.Seed
	nextMove.RoundNo = move.RoundNo + 1
	nextMove.OldBoard = move.NewBoard
	return nextMove
}

//TODO
func (move *moveT) InitFirstMove() {
	move.initMoveCollections()
	return
}

func marshalExcludeFields(structa interface{}, excludeFields map[string]bool,
) (jsona []byte, status error) {
	value := reflect.ValueOf(structa)
	typa := reflect.TypeOf(structa)
	size := value.NumField()
	jsona = append(jsona, '{')
	for i := 0; i < size; i++ {
		structValue := value.Field(i)
		var fieldName string = typa.Field(i).Name
		if marshalledField, marshalStatus := json.Marshal(
			(structValue).Interface()); marshalStatus != nil {
			return []byte{}, marshalStatus
		} else {
			if excludeFields[fieldName] == false {
				jsona = append(jsona, '"')
				jsona = append(jsona, []byte(fieldName)...)
				jsona = append(jsona, '"')
				jsona = append(jsona, ':')
				jsona = append(jsona, (marshalledField)...)
				if i != size-len(excludeFields) {
					jsona = append(jsona, ',')
				}
			}
		}
	}
	jsona = append(jsona, '}')
	return
}

func (move *moveT) initMoveCollections() {
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
}

// Get Seed
func (move *moveT) GetSeed() string {
	stringa := hex.EncodeToString(move.Seed.Bytes())
	return stringa
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
			move.NonMovedTiles = append(move.NonMovedTiles,
				ism.mergeHopefulIndex)
		} else {
			move.NonMergeMoves = append(move.NonMergeMoves,
				nonMergeMoveT{ism.mergeHopefulIndex,
					ism.mergeHopefulDestination})
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

func (move *moveT) randInt(upperLimit int, previous bool) int {
	hash := sha256.New()
	var roundNumberCast big.Int
	var upperLimitCast big.Int
	var bigInt big.Int
	roundNumberCast.SetInt64(int64(move.RoundNo))
	if previous {
		roundNumberCast.Add(&roundNumberCast, bigInt.SetInt64(-1))
	}
	upperLimitCast.SetInt64(int64(upperLimit))
	bigInt.Set(&move.Seed)
	bigInt.Add(&bigInt, &roundNumberCast)
	_, _ = hash.Write(bigInt.Bytes())
	var throwaway []byte
	bigInt.SetBytes(hash.Sum(throwaway))
	return int(bigInt.Rem(&bigInt, &upperLimitCast).Int64())
}

func (move *moveT) firstPass() {
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
		ism.currentIndex.nextPosition(ism.bigStep, 1)
	}
}

func (move *moveT) secondPass() {
	var mergePossibleColumns bool
	var lastValueColumns int
	var mergePossibleRows bool
	var lastValueRows int
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			rowPosition := positionT{i, j}
			columnPosition := positionT{j, i}
			rowElem := move.NewBoard.get(rowPosition)
			columnElem := move.NewBoard.get(columnPosition)
			if rowElem == 0 {
				move.NewTileCandidates = append(move.NewTileCandidates,
					rowPosition)
			} else {
				if lastValueRows == rowElem {
					mergePossibleRows = true
				}
			}
			if columnElem != 0 && columnElem == lastValueColumns {
				mergePossibleColumns = true
			}
			lastValueRows = rowElem
			lastValueColumns = columnElem
		}
	}
	var hasNewTiles bool
	if len(move.NewTileCandidates) == 0 {
		hasNewTiles = false
	} else {
		hasNewTiles = true
	}
	gameOn := mergePossibleRows || mergePossibleColumns || hasNewTiles
	move.IsGameOver = !gameOn

}

func (move *moveT) ResolveMove() {
	move.firstPass()
	move.secondPass()
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
