package mockDB

import (
	"testing"
)

func TestInitGameFull(t *testing.T) {
	qbm := QueryBuilderMockup{}
	qbm.ResetState()
	qbm.AddUser("adam")
	qb := QueryBuilder{}
	q := qb.InitGame("adam")
	q.Execute()
	var queryWorked bool
	queryWorked = <-*q.StatusChannel()
	if !queryWorked {
		t.Error("Unexpected QueryError")
	}
	if q.ResultString() != "it worked" {
		t.Error("Unexpected QueryResult")
	}
}
