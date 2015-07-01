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

import "../../queryBuilder"
import "fmt"

func main() {
	
	queryBuilder.Init()
	length := queryBuilder.ListLen()
	fmt.Println(length)
	number, roundNo, direction := queryBuilder.ReadJson()
	fmt.Println("1: ", number)
	fmt.Println("2: ", roundNo)
	fmt.Println("3: ", direction)
}
