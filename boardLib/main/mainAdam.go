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
import "bufio"
import "os"
import "log"
import "strings"
import "encoding/json"
import "../../API2048"

func prettyPrint(data string, oldBoard bool) {
	accessString := "NewBoard"
	if oldBoard {
		accessString = "OldBoard"
	}
	var result map[string]interface{}
	_ = json.Unmarshal([]byte(data), &result)
	parsedResult := result[accessString].([]interface{})
	fmt.Println("\n")
	for i := 0; i < API2048.BoardSize; i++ {
		fmt.Print("\t")
		parsedResultInner := parsedResult[i].([]interface{})
		for j := 0; j < API2048.BoardSize; j++ {
			fmt.Print(parsedResultInner[j], "\t")
		}
		fmt.Println("\n")
	}
	fmt.Println("Score ", result["RoundNo"])
	return
}

func main() {
	move := boardLib.CreateMove()
	move.InitFirstMove()
	prettyPrint(move.ExternalView(), true)
	var direction string
	var line string
	var err error
	for true {
		fmt.Println("press w, s, a or d and press enter.")
		reader := bufio.NewReader(os.Stdin)
		line, err = reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from stdin")
			os.Exit(1)
		}
		switch {
		case strings.Contains(line, "a"):
			direction = "left"
		case strings.Contains(line, "s"):
			direction = "down"
		case strings.Contains(line, "w"):
			direction = "up"
		case strings.Contains(line, "d"):
			direction = "right"
		default:
			continue
		}
		move.SetDirection(direction)
		move.ResolveMove()
		prettyPrint(move.ExternalView(), false)
		if move.GetGameOver() {
			break
		} else {
			move = move.CreateNextMove()
		}
	}
	fmt.Println("Game Over")
}
