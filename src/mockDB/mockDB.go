package mockDB

import (
	"API2048"
	"time"
)

type QueryError struct {
}

func (*QueryError) Error() string { return "" }

func (*QueryError) Code() int { return API2048.NoGameError }

type Query struct {
	Status chan bool
	Result string
}

func (q *Query) init() { q.Status = make(chan bool) }

func (*Query) Error() *QueryError { return nil }

func (q *Query) ResultString() string { return q.Result }

func (q *Query) Execute() {
	time.Sleep(time.Millisecond * 3)
	q.Result = "it worked"
	go func() { q.Status <- true }()
}

func (q *Query) StatusChannel() *chan bool {
	return &q.Status
}

type QueryBuilderMockup struct {
}

func (*QueryBuilderMockup) AddUser(userID string) {}

func (*QueryBuilderMockup) ResetState() {}

type QueryBuilder struct {
}

func (*QueryBuilder) InitGame(userID string) Query {
	q := Query{}
	q.init()
	return q
}
