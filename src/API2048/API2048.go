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

package API2048

// The size of the board. Most likely will always stay 4,
// as this is what seems to be most playable.
const BoardSize = 4

type QueryCallback interface {
	Callback(query Query)
}

// Query that is passed to the database (or in-memory database mock-up).
type Query interface {

	// Synchronous call, which waits at most maxMillis and by the time it is
	// finished, the Result or the Error will be ready. At most one call to
	// ExecuteAndWait can be made at one time.
	ExecuteAndWait(maxMillis int)

	// Asynchronous call. Check IsReady to see if it finished, or listen on the StatusChannel.
	Execute()

	// Asynchronous call with a callback. Pass in an implementation of
	// QueryCallback and Query.Callback will be called with the current instance
	// of this query when the result is ready.
	ExecuteAndCall(callback QueryCallback)

	// Channel, which receives true if a call to Execute was successful or false, if it errored.
	StatusChannel() chan bool

	// Result as a json string. Depending on the kind of Query it can be parsed
	// as a Move, a slice of ints or something else.
	ResultString() string

	// Returs the query string.
	QueryString() string

	// The most recent error or nil
	Error() QueryError

	// True if the call finished. If true, either Error or Result are guaranteed
	// to be available.
	IsReady() bool

	// If ready, how long has the query taken. How long since it was issued
	// otherwise. Result in Millis.
	HowLong() int
}

const (
	BackendConnectionError = iota
	NoGameError            = iota
)

type QueryError interface {
	Error() string
	Code() int
}

// QueryBuilder builds queries that can be subsequently executed.
type QueryBuilder interface {
	HasCurrentGame(userID string)

	// Prepares a query, which checks if there is a game already open for the
	// user. If there is it errors. If there is not, it
	// creates the first move. Then it updates the user list with the
	// current game. The result is an empty string.
	InitCurrentGame(userID string) Query

	//
	//
	CloseGame(userID string) Query

	// Prepares a query, which returns the most recent move. If there's no
	// current game it errors.
	// The result can be parsed as move.
	GetCurrentGameLastMove(userID string) Query

	// Prepares a query, which returns the last move of the current game of
	// the user. If there is no open game, it errors with the NoGame code.
	GetCurrentGame(userID string) Query

	// Prepares a query, which returns the score of [left, right] best games
	// played by the user so far. The range is inclusive, so to get the best and
	// the second best game one would use GetBestGameWindow(userID, 1, 2). If
	// the the interval is too broad, the query will return the broades result
	// possible. The result is parsable as a slice of GameScores.
	GetBestGames(userID string, left int, right int) Query

	// Prepares a query, which appends a move to the current game.
	AppendMoveCurrentGame(move string) Query
}

type QueryBuilderMockup interface {

	// Reinitalises the detabase mockup so that it is like new, uninitialised database.
	ResetState()

	// Adds a new user to the database mockup so games can be registered against that user.
	AddUser(userID string)
}

type MoveCreator interface {
	// Creates an empty, unitialised, unresolved Move.
	CreateMove() Move
}

// A single move representation. Contains all the logic necessary for resolving
// moves, parsing them from the database and serialising them in a form
// acceptable by the database. Implemented by boardLib.
type Move interface {
	// This initialises all the fields with sensible defaults for the first
	// move. This includes 2 random tiles, a random seed, a non-zero round
	// number, a valid direction and non-nil pointers for all internal data
	// structures. After initialisation it is safe to call other functions on
	// this struct.
	InitFirstMove()

	// Like InitFirstMove, but with more control over the initial conditions of
	// the move
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
	// validation. Resolved moves return isResolved true.
	// TODO refactor to include just the first error.
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
