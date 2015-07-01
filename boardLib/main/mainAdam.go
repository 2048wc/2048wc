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

import "../../boardLib"
import "fmt"
func main() {
	move := boardLib.CreateMove(
		[boardLib.BoardSize][boardLib.BoardSize]int{
			{16, 8, 4, 2},
			{4, 2, 2, 0},
			{2, 4, 0, 2},
			{2, 0, 0, 0},
		},
		"down",
		20,
		"e9ccc20fdb924ed423ad1b46c6df43516685f4c2bc36e202ad467af1b1d1febf",
	)
	boardLib.PrintBoard(move.OldBoard)
	move.ExecuteMove()
	fmt.Println(move)
	boardLib.PrintBoard(move.NewBoard)
}
