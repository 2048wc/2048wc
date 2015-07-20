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

package main

import "fmt"
import "../../queryBuilder"
import "github.com/fzzy/radix/redis"

func main() {
	ugqb := queryBuilder.QueryBuilder{}
	var query []string
	query = ugqb.InitGame("a01de45a674ad007a9abcf57b094a2a3", `{"move":true}`)
	query = query
	client, err := redis.Dial("tcp", "localhost:6379")
	replya, errora := client.Cmd("EVAL",
	   `redis.call("FLUSHALL")
        redis.call("RPUSH", KEYS[1], "prevGameid");
        redis.call("RPUSH", "prevGameid", "finish");
        local lastGameid = redis.call("LRANGE", KEYS[1], -1, -1); 
		if lastGameid[1] == nil then
			error("Failed to find the last gameid of this user"); 
		    return false; 
		end 
		local lastGameMove = redis.call("LRANGE", lastGameid[1], -1, -1); 
		if lastGameMove[1] == nil then 
			error("Failed to find the last move of the last game of this user") 
		    return false; 
		end 
		if lastGameMove[1] == "finish" then 
    		local call1 = redis.call("RPUSH", KEYS[1], KEYS[3]); 
    		local call2 = redis.call("RPUSH", KEYS[3], KEYS[2]);
            local checkuserid = redis.call("LRANGE", KEYS[1], -1, -1);
            local checkgameMove = redis.call("LRANGE", KEYS[3], -1, -1);
            print("call1", call1);
            print("call2", call2);
            print("checkuserid", checkuserid[1]);
            print("checkGameMove", checkgameMove[1]); 
			return true; 
		else 
		    return false; 
		end`,
		3, "a01de45a674ad007a9abcf57b094a2a3", `{"move":true}`, `bleh`).Int()
	fmt.Println(err, replya, "asd", errora)
	err = err
}

/*	*/

/*import "fmt"
import "github.com/fzzy/radix/redis"

func main() {
    client, err := redis.Dial("tcp", "localhost:6379")
    if err != nil {
        println("error! could not coonect to db")
    } else {
            client.Cmd("SET", "foo", "bar")
            reply, err2 := client.Cmd("EVAL", "return 1", 0).Str()
            if err2 != nil {
            println("errora: ", err2)
        }
        fmt.Println(reply)
    }
}
*/
