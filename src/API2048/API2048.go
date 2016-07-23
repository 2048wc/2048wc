package API2048

/* Copyright (C) 2015, 2016  Adam Kurkiewicz and Iva Babukova
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

/********************************** Move API **********************************/

// The size of the board. Most likely will always stay 4,
// as this is what seems to be most playable.
const BoardSize = 4

// Any component implementing logic of 2048 must satisfy this interface.
// Implemented by boardLib.
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

/*********************************** DB API ***********************************/

// A move representation that can be serialised into a database. Move field
// contains an object implementing the Move interface encoded as a json
// string, with some fields exposed in unserialised form (RoundNo -> MoveNo,
// GameOver -> !IsOpen). The purpose of this indirection is to avoid serialising
//  Move into database directly in order to keeping the implementations of
// DB API and Move API independent.
type MoveEntry struct {
	// json-encoded Move
	Move string

	// Same as RoundNo in Move
	MoveNo int

	// Inverse of GameOver in Move
	IsOpen bool
}

type ScoreRange struct {
	UserID string
	Scores []int
}

type QueryCallback interface {
	Callback(query Query)
}

// Query that is passed to the database (or in-memory database mock-up).
type Query interface {

	// Synchronous call, which waits at most maxMillis and by the time it is
	// finished, the Result or the Error will be ready. If you pass 0, it will
	// wait 20s (20000ms).
	ExecuteAndWait(maxMillis int)

	// Asynchronous call. Check IsReady to see if it finished, or listen on the StatusChannel.
	Execute()

	// Asynchronous call with a callback. Pass in an implementation of
	// QueryCallback and Query.Callback will be called with the current instance
	// of this query when the result is ready.
	ExecuteAndCall(callback QueryCallback)

	// Channel, which receives true if a call to Execute was successful or false, if it errored.
	StatusChannel() *chan bool

	// Type depends on the kind of query. Some queries will return MoveEntry, others will return a range of MoveEntries.
	Result() interface{}

	// Returs the query string.
	// QueryString() string

	// The most recent error or nil
	Error() QueryError

	// True if the call finished. If true, either Error or Result are guaranteed
	// to be available.
	IsReady() bool

	// If ready, how long has the query taken. How long since it was issued
	// otherwise. Result in Millis.
	HowLong() int
}

// DB API error codes.
const (
	BackendConnectionError = iota
	NoGameError            = iota
	NoUserError            = iota
	MoveConstraintError    = iota
)

// DB API error. TODO refactor into a struct
type QueryError interface {
	Error() string
	Code() int
}

// QueryBuilder builds queries that can be subsequently executed.
//
// Game can be either:
// - inexistent
// - closed
// - open
// The normal flow is inexistent -> open -> closed
//
// Any call in this file can return NoUserError or BackendConnectionError
type QueryBuilder interface {

	// Saves MoveEntry to a database if certain constraints are met.
	//
	// Either:
	// 1. the previous move is open and
	// 2. the previous move's MoveNo is one less than the current move's MoveNo
	// Or:
	// 3. the previous move is closed or inexistent and
	// 4. the current move's MoveNo is 0.
	//
	// These constraints are called MoveConstraints, if you don't fulfill them
	// you get MoveConstraintError
	//
	// Other error you can get is NoUserError
	//
	// returnedQuery.Result() can be safely cast to MoveEntry
	SetNextMove(userID string, dbentry MoveEntry) (returnedQuery Query)

	// Prepares a query, which returns the most recent move. If there's no
	// current game it errors with NoGameError.
	//
	// Query.Result() is of type MoveEntry
	GetLastMove(userID string) Query

	// Prepares a query, which returns the score of [left, right] best games
	// played by the user so far. The range is inclusive, so to get the best and
	// the second best game one would use GetBestGames(userID, 1, 2). If
	// the the interval is too broad, the query will return the broadest result
	// possible. The result is parsable as a slice of GameScores.
	//GetBestGames(userID string, left int, right int) Query
}

type QueryBuilderMockup interface {

	// Reinitalises the detabase mockup so that it is like new, uninitialised database.
	ResetState()

	// Adds a new user to the database mockup so games can be registered against that user.
	ResetUser(userID string)
}
