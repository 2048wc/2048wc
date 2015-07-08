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

//import "../../boardLib"
import "fmt"
import "reflect"
import "encoding/json"

func MarshalOnlyFields(structa interface{},
	includeFields map[string]bool) (jsona []byte, status error) {
	value := reflect.ValueOf(structa)
	typa := reflect.TypeOf(structa)
	size := value.NumField()
	jsona = append(jsona, '{')
	for i := 0; i < size; i++ {
		structValue := value.Field(i)
		var fieldName string = typa.Field(i).Name
		if marshalledField, marshalStatus := json.Marshal((structValue).Interface()); marshalStatus != nil {
			return []byte{}, marshalStatus
		} else {
			if includeFields[fieldName] {
				jsona = append(jsona, '"')
				jsona = append(jsona, []byte(fieldName)...)
				jsona = append(jsona, '"')
				jsona = append(jsona, ':')
				jsona = append(jsona, (marshalledField)...)
				if i+1 != len(includeFields) {
					jsona = append(jsona, ',')
				}
			}
		}
	}
	jsona = append(jsona, '}')
	return
}

type magic struct {
	Magic1 int
	Magic2 string
	Magic3 [2]int
}

func main() {
	var magic = magic{0, "tusia", [2]int{0, 1}}
	if json, status := MarshalOnlyFields(magic, map[string]bool{"Magic1": true}); status != nil {

	} else {
		fmt.Println(string(json))
	}

}
