package mockDB

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

import (
	"API2048"
	"fmt"
	"sync"
	"time"
)

//////////////////// MoveEntry ////////////////////

type MoveEntry API2048.MoveEntry

func (me *MoveEntry) InitMoveEntry() {
	me.IsOpen = true
}

////////////////////  Global State ////////////////////

type game []API2048.MoveEntry

type db map[string][]game

var dbLock sync.Mutex

var database db

////////////////////  QueryBuilderMockup ////////////////////

type QueryBuilderMockup struct {
}

func (*QueryBuilderMockup) ResetUser(userID string) {
	database[userID] = make([]game, 0, 10)
}

func (*QueryBuilderMockup) ResetState() {
	database = make(db)
}

////////////////////  QueryError ////////////////////

type QueryError struct {
	err  string
	code int
}

func (e QueryError) Error() string { return e.err }

func (e QueryError) Code() int { return e.code }

////////////////////  Query ////////////////////

type Query struct {
	Status    chan bool
	result    API2048.MoveEntry
	ToExecute func()
	err       QueryError
	ready     bool
}

func (q *Query) init() { q.Status = make(chan bool) }

func (q *Query) Error() API2048.QueryError {
	return q.err
}

func (q *Query) Result() interface{} {
	return q.result
}

func (q *Query) Execute() {
	time.Sleep(time.Millisecond * 3)
	q.ToExecute()
}

func (q *Query) ExecuteAndCall(callback API2048.QueryCallback) {
	return
}

func (q *Query) HowLong() int {
	return 0
}

func (q *Query) IsReady() bool {
	return q.ready
}

func (q *Query) ExecuteAndWait(millis int) {
	return
}

func (q *Query) StatusChannel() *chan bool {
	return &q.Status
}

////////////////////  QueryBuilder ////////////////////

type QueryBuilder struct{}

func (*QueryBuilder) SetNextMove(userID string, moveEntry API2048.MoveEntry) API2048.Query {
	q := &Query{}

	noUserMsg := "User " + userID + " doesn't exist"

	toExecuteTerminate := "toExecuteFinished"

	q.init()

	lockDB := func() {
		dbLock.Lock()
	}

	unlockDB := func() {
		dbLock.Unlock()
	}

	fail := func(err QueryError) {
		q.ready = true
		q.err = err
		go func() { q.Status <- false }()
	}

	succeed := func() {
		q.ready = true
		q.result = moveEntry
		go func() { q.Status <- true }()
	}

	existsLastGame := func() bool {
		games, userExists := database[userID]
		if !userExists {
			err := QueryError{noUserMsg, API2048.NoUserError}
			defer fail(err)
			panic(toExecuteTerminate)
		}
		return len(games) != 0
	}

	getLastGame := func() (int, game) {
		games, userExists := database[userID]
		if !userExists {
			err := QueryError{noUserMsg, API2048.NoUserError}
			defer fail(err)
			panic(toExecuteTerminate)
		}
		nextGame := len(games)
		lastGame := nextGame - 1
		// This line of code normalises Case 1.
		var g game
		if lastGame < 0 {
			lastGame = 0
			g = nil
		} else {
			g = games[lastGame]
		}
		return lastGame, g
	}

	getLastMove := func(g game) (int, API2048.MoveEntry) {
		nextMove := len(g)
		lastMove := nextMove - 1
		var m API2048.MoveEntry
		if lastMove < 0 {
			lastMove = 0
			m = API2048.MoveEntry(MoveEntry{})
		}

		return lastMove, m
	}

	extendLastGame := func(me API2048.MoveEntry) {
		_, lastGame := getLastGame()
		previousMoveIndex, previousMoveEntry := getLastMove(lastGame)
		cond1 := previousMoveEntry.MoveNo+1 == me.MoveNo
		cond2 := previousMoveEntry.IsOpen
		if cond1 && cond2 {
			lastGame[previousMoveIndex+1] = me
			defer succeed()
			panic(toExecuteTerminate)
		} else {
			defer fail(QueryError{"previous and next move's number don't match or previous move is not open", API2048.MoveConstraintError})
			panic(toExecuteTerminate)
		}
	}

	beginNewGame := func(me API2048.MoveEntry) {
		_, lastGame := getLastGame()
		_, lastMove := getLastMove(lastGame)
		cond1 := me.MoveNo == 0
		cond2 := len(lastGame) == 0
		var cond3 bool // = false
		if !cond2 {
			cond3 = !lastMove.IsOpen
		}
		if cond1 && (cond2 || cond3) {
			lastGame = append(lastGame, moveEntry)
			defer succeed()
			panic(toExecuteTerminate)
		} else {
			defer fail(QueryError{"you're trying to create a new game and your me.MoveNumber is not 0, or your previous move is neither closed nor inexistent", API2048.MoveConstraintError})
			panic(toExecuteTerminate)
		}
	}
	q.ToExecute = func() {
		defer func() {
			r := recover()
			if r != toExecuteTerminate {
				print("unexpected panic:", r)
				panic(r)
			}
		}()
		func() {
			// There are three posibilities, really. Imagine games are boxes. Each box can contain a sequence of open moves O, closed moves X or no moves -
			//
			// Case 1 (no last game)
			// lastGame ------
			//
			// Case 2 (open last game)
			// lastGame 0000--
			//
			// Case 3 (closed last game)
			// lastGame 00X---
			lockDB()
			defer unlockDB()
			case1 := !existsLastGame()
			fmt.Println("")
			if case1 {
				beginNewGame(moveEntry)
				return
			}
			_, lastGame := getLastGame()
			_, lastMove := getLastMove(lastGame)
			case2 := lastMove.IsOpen
			case3 := !lastMove.IsOpen

			if case2 {
				extendLastGame(moveEntry)
			}

			if case3 {
				beginNewGame(moveEntry)
			}
		}()

	}
	return q
}

func (*QueryBuilder) GetLastMove(userID string) (returnedQuery API2048.Query) {
	q := &Query{}
	q.init()
	return q
}
