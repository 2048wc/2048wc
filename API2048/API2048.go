package API2048


// The size of the board. Most likely will always stay 4,
// as this is what seems to be most playable.
const BoardSize = 4

// Query that is passed to the database (or in-memory database mock-up).
type Query interface {
	
	// Synchronous call, which waits at most maxMillis and by the time it is
	// finished, the Result or the Error will be ready. 
	ExecuteAndWait(maxMillis int)
	
	// asynchronous call. Check IsReady to see if it finished.
	Execute()

	// Result as a json string. Depending on the kind of Query it can be parsed 
	// as a Move, a slice of ints or something else.
	GetResult() string
	
	// Returs the query string. 
	GetQuery() string

	// Any errors can be accessed here
	GetError() error
	
	// True if the call finished. If true, either Error or Result are guaranteed
	// to be available.
	IsReady() bool

	// If ready, how long has the query taken. How long since it was issued
	// otherwise. Result in Millis.
	HowLong() int
}

// QueryBuilder builds queries that can be subsequently executed.
type QueryBuilder interface {

	// prepares a query, which checks if there is a game already open for the
	// user. If there is it errors. If there is not, it
	// creates the first move. Then it updates the user list with the
	// current game. The result is an empty string.
	InitGame(userID string, firstMove string) Query

	// prepares a query, which returns the most recent move. If there's no
	// current game it errors.
	// The result can be parsed as move.
	GetCurrentGameLastMove(userID string) Query

	// prepares a query, which returns the last move of the current game of
	// the user. If there is no open game, it errors
	GetCurrentGame(userID string) Query

	// prepares a query, which returns the score of [left, right] best games
	// played by the user so far. The range is inclusive, so to get the best and
	// the second best game one would use GetBestGameWindow(userID, 1, 2). If
	// the the interval is too broad, the query will return the broades result
	// possible. The result is parsable as a slice of GameScores.
	GetBestGames(userID string, left int, right int) Query

	// Adds move json if the round number is 1 more than the last stored round
	// number. Error otherwise.
	AppendMove(move string, gameID string) Query
}


// A single move representation. Contains all the logic necessary for resolving
// moves, parsing them from the database and serialising them in a form
// acceptable by the database.
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

