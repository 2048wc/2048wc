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
import "../boardLib"

func Init() {
	move := boardLib.CreateMove(boardLib.BoardT{
		{0,0,0,1},
		{0,2,3,4}, 
		{0,7,6,4},
		{1,1,2,2}},
		"left", 25, "asdfasdfasdf")
	move.ExecuteMove()
	jsonified, _ := json.Marshal(move)
	fmt.Println(string(jsonified))
	client, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		// handle err
		fmt.Println("in init: can't connect")
	}

	client.Cmd("DEL", "queryBuilder::Iva").Str()
	client.Cmd("RPUSH", "queryBuilder::Iva", `{"RoundNo":17, "direction":"left"}`).Str()
	client.Cmd("RPUSH", "queryBuilder::Iva", "17").Int()
	client.Close()
}

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

func ReadJson() (int, int, string){
	
	client, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		// handle err
		fmt.Println("in dial: can't connect")
	}
	
	jsona, err := client.Cmd("LRANGE", "queryBuilder::Iva", -2, -1).List()
	if err != nil {
		fmt.Println("error in reading from the list")
	}
	fmt.Println("jsona: ", jsona[0])
	
	type RoundNoDirection struct {
		RoundNo int
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

func Foo() {
	
	/*
	friends, err := client.Cmd("SET", "friends", "iva").Str()
	if err != nil {
		fmt.Println("couldn't set friends")
	}
	foo, err := client.Cmd("GET", "friends").Str()
	if err != nil {
			// handle err
			fmt.Println("couldn't get friends")
	}

	fmt.Println(friends)
	fmt.Println(foo)
	*/
}
