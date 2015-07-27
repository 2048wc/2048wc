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
	fmt.Println(err)
	client = client
	/*
	replya, errora := client.Cmd("EVAL", query[0], 3, query[1], query[2], query[3]).Int()
	if errora != nil {
		println("there is error!!!!");
	}

	fmt.Println(replya, "asd", errora)
	*/
	err = err
}

