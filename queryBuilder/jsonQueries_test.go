package queryBuilder

import "testing"

func TestInterfaceFullfiled(t *testing.T) {
	var toFulfill UsersGamesQueryBuilder
	toFulfill = QueryBuilder{}
	toFulfill = toFulfill
	return
}


// USERID validation tests:
func TestGameIdValidationEmptyString(t *testing.T) {
	var gameid string = ""
	if isGameIdValid(gameid) {
		t.Error("gameIDValidation function failed fro empty string")
		t.Fail()
	}
	return
}

func TestGameIdValidationBiggerString(t *testing.T) {
	var gameid string = "aaadadadadadaddabcbcbcbcbdbdbabab2222323124431b23bb3222123232232233232423242342bb3b2bbabbcdbefedbbb"
	if isGameIdValid(gameid) {
		t.Error("gameIDValidation function failed fro empty string")
		t.Fail()
	}
	return
}

func TestGameIdValidationNotValidChars(t *testing.T) {
	var gameid string = "abcaaaaaaaaaaaaaaaaaaaaaaaaaaaaz"
	if isGameIdValid(gameid) {
		t.Error("gameIDValidation function failed fro empty string")
		t.Fail()
	}
	return
}

func TestGameIdValidationValidChars(t *testing.T) {
	var gameid string = "abcdeaaa674ad007a9abcf57aaa4a2a3"
	if !isGameIdValid(gameid) {
		t.Error("gameIDValidation function failed fro empty string")
		t.Fail()
	}
	return
}

func TestGameIdValidationInvalidChars(t *testing.T) {
	var gameid string = `123, -1, -1); redis.call("DELETE", "ALL") #`
	if isGameIdValid(gameid) {
		t.Error("gameIDValidation function failed")
		t.Fail()
	}
	return
}

// InitGame tests:
func TestInitGameInvalidUID(t *testing.T) {
	qb := &QueryBuilder{}
	query := qb.InitGame("ivababukova", "left")
	if (query != nil ) {
		t.Error("didn't fail when the input was invalid")
		t.Fail()
	}
}

func TestInitGameValidUID(t *testing.T) {
	qb := &QueryBuilder{}
	query := qb.InitGame("a01de45a674ad007a9abcf57b094a2a3", "left")
	if (query == nil ) {
		t.Error("failed although the input was valid")
		t.Fail()
	}
	
	if(query[0] != "lastGameID = eval \"return edis.call(\"LRANGE\", KEYS[1], -1, -1)\" 1 a01de45a674ad007a9abcf57b094a2a3\n") {
		t.Error("initGame doesn't return a correct query line 1 to initiate a game")
		t.Fail()
	}
	
	if(query[1] != "lastItemInGame = eval \"return redis.call(\"LRANGE\", KEYS[1], -1, -1)\" 1 lastGameID\n") {
		t.Error("initGame doesn't return a correct query line 2 to initiate a game")
		t.Fail()
	}
	
	if(query[2] != "if lastItemInGame == \"finished\" then\n") {
		t.Error("initGame doesn't return a correct query line 3 to initiate a game")
		t.Fail()
	}
	if(query[3] != "eval \"return redis.call('SET',KEYS[1],ivasgame)\" 1 a01de45a674ad007a9abcf57b094a2a3\n") {
		t.Error("initGame doesn't return a correct query line 4 to initiate a game")
		t.Fail()
	}
	
	if(query[4] != "eval \"return redis.call('SET',KEYS[1],left)\" 1 ivasgame\n") {
		t.Error("initGame doesn't return a correct query line 5 to initiate a game")
		t.Fail()
	}
	
	if(query[5] != "else\n") {
		t.Error("initGame doesn't return a correct query line 6 to initiate a game")
		t.Fail()
	}
	
	if(query[6] != "return \"kupa\"") {
		t.Error("initGame doesn't return a correct quety to initiate a game")
		t.Fail()
	}
}
