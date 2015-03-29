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

import "github.com/2048wc/2048wc/boardLib"
import "fmt"

func main() {
  oldBoard := boardLib.BoardT{
		[boardLib.BoardSize]int{16, 8, 4, 2},
		[boardLib.BoardSize]int{4, 2, 2, 0},
		[boardLib.BoardSize]int{2, 4, 0, 2},
		[boardLib.BoardSize]int{0, 0, 0, 0},
	}
	move := boardLib.CreateMove(
		oldBoard,
		"right",
                24,
		"e9ccc20fdb924ed423ad1b46c6df43516685f4c2bc36e202ad467af1b1d1febf",
	)
	fmt.Println(move.ComputeDistance())
	
}
