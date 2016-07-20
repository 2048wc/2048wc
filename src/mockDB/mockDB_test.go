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
	"reflect"
	"testing"
)

func TestInitGameFull(t *testing.T) {
	var qbm API2048.QueryBuilderMockup
	qbm = &QueryBuilderMockup{}
	qbm.ResetState()
	qbm.ResetUser("adam")
	var qb API2048.QueryBuilder
	qb = &QueryBuilder{}
	var me MoveEntry
	me = MoveEntry{"move1", 0, true}
	me.InitMoveEntry()
	q := qb.SetNextMove("adam", API2048.MoveEntry(me))
	q.Execute()
	var queryWorked bool
	queryWorked = <-*q.StatusChannel()
	if !queryWorked {
		t.Log(q.Error().Error())
		t.Error("Unexpected QueryError")
	}
	result := q.Result()
	expected := me
	if reflect.DeepEqual(result, expected) {
		t.Log(result)
		t.Log(expected)
		t.Error("Unexpected QueryResult")
	}
}

func TestInitGameNoUser(t *testing.T) {
	qbm := QueryBuilderMockup{}
	qbm.ResetState()
	qb := QueryBuilder{}
	var me MoveEntry
	me = MoveEntry{}
	me.InitMoveEntry()
	q := qb.SetNextMove("adam", API2048.MoveEntry(me))
	q.Execute()
	queryWorked := <-*q.StatusChannel()
	if queryWorked {
		t.Error("Expected Error")
	} else {
		if q.Error().Code() != API2048.NoUserError {
			t.Error("UnexcpectedError Code")
		}
	}
}
