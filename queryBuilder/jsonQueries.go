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

package queryBuilder

import "fmt"
import "github.com/fzzy/radix/redis"
import "strconv"
import "strings"
import "encoding/json"
import "log"
import "io"

//import "../boardLib"
import "time"
import "math"
import "errors"

// import "container/list"
import "reflect"

/* uncomment and use this once the api is fully implemented
type UserGameQuery struct {
	query []string
	replyParser *func(query *UserGameQuery) interface{}
	replyParserReturnType
}
*/

type UsersGamesQueryBuilder interface {

	// prepares a query, which checks if there is a game already open for the
	// user. If there is it kupas. If there is not, it
	// creates a move. Then it updates the user list with the current game.
	InitGame(userID string, initialMove string) []string

	// prepares a query, which, if there is a current game, it returns the most
	// recent move. Otherwise, kupa.
	GetCurrentGameLastMove(userID string) []string

	// prepares a query, which returns the last move of the current game of
	// the user. If there is no open game, it kupas ?calles InitGame?
	GetCurrentGame(userID string) []string

	// prepares a query, which returns the score of the best game played by the
	// user so far. If the user hasn't played any game, it kupas
	GetBestGameWindow(userID string, n int, m int) []string

	// Adds move json only if roundnumber in the json corresponds to
	// 1 + roundnumber stored as the last element in the games list
	// returns 0 if successful or -1 if not succcessful
	// AppendMoveToGame (jmove string, gameID string) error
}
/*
type DBQuery interface {
	toString() string
	configProperty(string, string)
	executeQuery() DBResponse
}

type DBError interface {
	webTierError bool // unrecoverable
	DBUnreachable bool // retry query 3 times and then assume DBBroke
	DBBroke bool // ask consul for a route to a different DB instance
	errora error
}

type DBResponse interface {
	isSuccess() bool
	errora DBError
}
*/
/*

	RetryDatabase(call *func() (bool, error), totalDelay time.Duration, numberOfTries int) error	
	// this function tries to call func *call* *numberOfTries* times for the
	// interval of *totalDelay time*. It returns an error message: for succcess
	// it returns "200 OK", if failed, it returns the message of the error


*/

func delayMiliseconds(n time.Duration) {
	time.Sleep(n * time.Millisecond)
}

func isGameIdValid(gameID string) bool {
	// gameID is 32 chars
	// gameID doesn't contain bad chars
	valid := map[byte]bool{'0': true, '1': true, '2': true, '3': true,
		'4': true, '5': true, '6': true, '7': true, '8': true, '9': true,
		'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true}
	const length = 32
	if len(gameID) != length {
		return false
	}
	for i := 0; i < length; i++ {
		if valid[gameID[i]] == false {
			return false
		}
	}
	return true
}

type QueryBuilder struct{}

func (qb QueryBuilder) GetBestGameWindow(userID string, n int, m int) []string {
	return []string{""}
}

// prepares a query, which checks if there is a game already open for the
// user. If there is it kupas. If there is not, it
// creates a move. Then it updates the user list with the current game.
func (qb QueryBuilder) InitGame(userID string, initialMove string) []string {

	// TODO make a initialMove validator?
	// TODO implement a gameID generator
	gameID := "ivasgame"
	gameID = gameID
	var query []string
	query = make([]string, 0, 0)
	query = append(query,
			 `local lastGameid = redis.call("LRANGE", KEYS[1], -1, -1); 
			 if lastGameid[1] == nil then error("Failed to find the last gameid of this user") 
		         return false; 
			 end 
			 local lastGameMove = redis.call("LRANGE", lastGameid[1], -1, -1); 
			 if lastGameMove[1] == nil then error("Failed to find the last move of the last game of this user") 
		         return false; 
			 end 
			 if lastGameMove[1] == "finish" then 
    		     redis.call("RPUSH", KEYS[1], KEYS[2]); 
    		     redis.call("RPUSH", KEYS[2], KEYS[3]); 
			     return true; 
			 else 
		         return false; 
			 end 
			 print(redis.call("LRANGE", lastGameid[1], -1. -1));`,
	 userID, initialMove, gameID)
	return query
}

func (qb QueryBuilder) GetCurrentGameLastMove(userID string) []string {
	var query []string
	query[0] =
			 "local lastGameid = redis.call(\"LRANGE\", KEYS[1], -1, -1);" +
			 "if lastGameid[1] == nil then" +
			 "    error(\"Failed to find the last gameid of this user\");" +
		     "    return \"\";" +
			 "end" +
			 "local lastGameMove = redis.call(\"LRANGE\", lastGameid[1], -1, -1);" +
			 "if lastGameMove[1] == nil then" +
			 "    error(\"Failed to find the last move of the last game of this user\")" +
		     "    return \"\"" +
			 "end" +
			 "if lastGameMove[1] == \"finish\" then" +
			 "    error(\"Error: the user has finished all his games\")" +
		     "    return \"\";" +
			 "end" +
			 "else" +
			 "    return lastGameMove[1]"
	query[1] = userID
	return []string{""}
}

func (qb QueryBuilder) GetCurrentGame(userID string) []string {
	return []string{""}
}

func (qb QueryBuilder) GetBestGame(userID string) []string {
	return []string{""}
}

// TODO finish the function implementation
func (qb QueryBuilder) AddMoveToDB(jmove string, listName string) error {

	err, client := dialDb()
	if err != nil {
		return err
	}
	prev_RoundNo, roundNo, _ := returnMoveFromDb(client)
	if prev_RoundNo+1 != roundNo {
		err := errors.New("Error: could not write move to the database")
		return err
	}

	// TODO write jmove to the database
	// TODO write the roundNumber to the database
	return errors.New("200 OK")
}

func (qb QueryBuilder) RetryDatabase(call *func() (bool, error), totalDelay time.Duration, numberOfTries int) error {

	if totalDelay <= 0 {
		err := errors.New("total delay is equal or less than 0")
		return err
	}
	totalDelay = totalDelay - time.Duration(numberOfTries*20)
	x := int(totalDelay) / int(math.Pow(2, float64(numberOfTries))-1)

	var success bool
	var err error
	if x <= 0 {
		err1 := errors.New("total delay given is too small")
		return err1
	}
	for i := 0; i < numberOfTries; i++ {
		success, err = (*call)()
		if success == false {
			time.Sleep(time.Duration(x) * time.Millisecond)
			x = x * 2
		} else if success {
			return err
		}
	}
	return err
}

/*
This function establishes a connection to the database.
If the connection was successfuly established, it returns 200 OK and pointer to
redis.client obtained.
If the connection was not established, it returns the error generated and nil
*/
func dialDb() (error, *redis.Client) {
	var err error
	client, err := redis.Dial("tcp", "localhost:6379")
 	if err != nil {
		return err, client
	} else {
		err = nil
		returnMoveFromDb(client)
	}
	return err, client

}

/* This function is called from RetryDatabase. If the connection is successfull,
it returns 200 OK and calls ReturnMoveFromDb. If there is an error, it returns the error message
*/
func ReadFromDb() (bool, error) {
	var err error
	var success bool
	err, client := dialDb()
	if err != nil {
		return success, err
	} else {
		err = errors.New("200 OK")
		fmt.Println("client: ", client)
		returnMoveFromDb(client)
		success = true
	}
	return success, err
}

// TODO return a whole Move, write a small fn to return the query
/* returns number, roundNo, direction from the move recorded in the database or -1
for all ints ans nil for all strings if the connection was not successful
*/
func returnMoveFromDb(client *redis.Client) (int, int, string) {
	jsona, err := PrintElementsFromDb(client, "queryBuilder::Iva", -2, -1)
	fmt.Println("type of jsona: ", reflect.TypeOf(jsona))
	if err != nil {
		fmt.Println("could not get the elements of te list requested")
		return -1, -1, ""
	}
	fmt.Println("jsona: ", jsona[0])

	type RoundNoDirection struct {
		RoundNo   int
		Direction string
	}

	var rd RoundNoDirection
	dec := json.NewDecoder(strings.NewReader(jsona[0]))
	for {
		if err := dec.Decode(&rd); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Println("roundNo: ", rd.RoundNo, "direction: ", rd.Direction)
	}

	number, _ := strconv.ParseInt(jsona[1], 10, 32)
	fmt.Println(number)

	client.Close()
	return int(number), rd.RoundNo, rd.Direction
}

func CreateLuaScript(jmove string) string {
	var luaScript string
	var moveS []string
	fmt.Println(jmove)
	moveS = strings.Split(jmove, ",")
	fmt.Println("splited jmove: ", moveS[1])
	// TODO write the lua script that has to be returned
	return luaScript
}

/* Functions that query the database */
func PrintElementsFromDb(client *redis.Client, list string, startIndex int, endIndex int) ([]string, error) {
	lista, err := client.Cmd("LRANGE", list, startIndex, endIndex).List()
	return lista, err
}

/* pushes specified string element in the specified list in the database */
// TODO make it accept elements of generic type
func PushElementToDb(client *redis.Client, list string, element string) {
	fmt.Println(client.Cmd("RPUSH", list, element))
	return
}

// TODO this function should return error, it is supposed to be called by RetryDatabase
// returns the current length of the list:
func ListLen() int {

	client, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		// handle err
		fmt.Println("in dial: can't connect")
	}

	length, err := client.Cmd("LLEN", "queryBuilder::Iva").Int()
	if err != nil {
		// handle err
		fmt.Println("can't check the lenght of queryBuilder::Iva")
	}

	client.Close()
	return length
}
