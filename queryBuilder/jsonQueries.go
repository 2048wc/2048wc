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

/*
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
	GetCurrentGameLastMove() []string

	// prepares a query, which returns the last move of the current game of
	// the user. If there is no open game, it kupas ?calles InitGame?
	GetCurrentGame(userID string) []string

	// prepares a query, which returns the score of the best game played by the
	// user so far. If the user hasn't played any game, it kupas
	GetBestGameWindow(userID string, n int, m int) []string

	// this function tries to call func *call* *numberOfTries* times for the
	// interval of *totalDelay time*. It returns an error message: for succcess
	// it returns "200 OK", if failed, it returns the message of the error
	RetryDatabase(call *func() (bool, error), totalDelay time.Duration, numberOfTries int) error

	// Adds move json only if roundnumber in the json corresponds to
	// 1 + roundnumber stored as the last element in the games list
	// returns 0 if successful or -1 if not succcessful
	// AppendMoveToGame (jmove string, gameID string) error
}

func delayMiliseconds(n time.Duration) {
	time.Sleep(n * time.Millisecond)
}

type QueryBuilder struct{}

// prepares a query, which checks if there is a game already open for the
// user. If there is it kupas. If there is not, it
// creates a move. Then it updates the user list with the current game.
func InitGame(userID string, initialMove string) string {

	var query string
	query = "lastGameID = redis.call(\"LRANGE\","+ userID + ", -1, -1)\n" +
	"lastItemInGame = redis.call(\"LRANGE\", lastGameID, -1, -1)\n" +
	"if lastItemInGame == \"finished\" then\n" +
	// init the game using initial move
	// TODO compose a gameid
    "else\n" + 
    "return \"kupa\"\n"
	
	fmt.Println(query, "\n")
	
	// fmt.Println(fmt.Sprintf("<%s>", userID))
	// lastGameID, err := client.Cmd("EVAL", fmt.Sprintf(`return redis.call("LRANGE",KEYS[1],"%s")`, "ivaLikes"), 1, "adam")

	/*
			client.Cmd("DEL", "queryBuilder::Iva").Str()
			client.Cmd("RPUSH", "queryBuilder::Iva", jsonified)
			client.Cmd("RPUSH", "queryBuilder::Iva", "17").Int()
			err = errors.New("200 OK")
		}
		return success, err*/

	return query
}

func (qb QueryBuilder) GetCurrentGameLastMove() []string {
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

func RetryDatabase(call *func() (bool, error), totalDelay time.Duration, numberOfTries int) error {

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
	var client *redis.Client = nil
	client, err = redis.Dial("tcp", "localhost:6379")
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
