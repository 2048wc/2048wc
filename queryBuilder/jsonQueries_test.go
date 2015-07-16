package queryBuilder

import "testing"

func TestInterfaceFullfiled(t *testing.T) {
	var toFulfill UsersGamesQueryBuilder
	toFulfill = QueryBuilder{}
	toFulfill = toFulfill
	return
}

func TestGameIdValidationEmptyString(t *testing.T) {
	var gameid string = ""
	if IsGameIdValid(gameid) {
		t.Error("gameIDValidation function failed fro empty string")
		t.Fail()
	}
	return
}

func TestGameIdValidationBiggerString(t *testing.T) {
	var gameid string = "aaadadadadadaddabcbcbcbcbdbdbabab2222323124431b23bb3222123232232233232423242342bb3b2bbabbcdbefedbbb"
	if IsGameIdValid(gameid) {
		t.Error("gameIDValidation function failed fro empty string")
		t.Fail()
	}
	return
}

func TestGameIdValidationNotValidChars(t *testing.T) {
	var gameid string = "abcaaaaaaaaaaaaaaaaaaaaaaaaaaaaz"
	if IsGameIdValid(gameid) {
		t.Error("gameIDValidation function failed fro empty string")
		t.Fail()
	}
	return
}

func TestGameIdValidationValidChars(t *testing.T) {
	var gameid string = "abcdeaaa674ad007a9abcf57aaa4a2a3"
	if !IsGameIdValid(gameid) {
		t.Error("gameIDValidation function failed fro empty string")
		t.Fail()
	}
	return
}

func TestGameIdValidationInvalidChars(t *testing.T) { // iva spake
	var gameid string = `123, -1, -1); redis.call("DELETE", "ALL") #`
	if IsGameIdValid(gameid) {
		t.Error("gameIDValidation function failed fro empty string")
		t.Fail()
	}
	return
}
