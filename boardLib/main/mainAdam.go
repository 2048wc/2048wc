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
import "math/big"
import "encoding/hex"

func main() {
	//128 64 32 16 8 4 2 1
	//1   0  1  1  1 1 0 1
	//1 + 4 + 8 + 16 + 32 + 128 = 189
	slica := []byte("\xbd\xbd")
	fmt.Println(len(slica))
	bigInta := (&big.Int{}).SetBytes(slica)
	stringa := hex.EncodeToString(bigInta.Bytes())
	fmt.Println(stringa)
}
