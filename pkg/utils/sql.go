//
// Copyright 2023 Beijing Volcano Engine Technology Ltd.
// Copyright 2023 Guangzhou Laboratory
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"strings"
	"unicode"

	"gorm.io/gorm"
)

// SearchWordFilter splits keyword by special characters, and the relation
// between keywords is 'OR'. In addition, we don't need to escape character
// because special characters has been removed as delimiters.
func SearchWordFilter(db *gorm.DB, word string, fields []string, exact bool) *gorm.DB {
	if len(fields) == 0 {
		return db
	}
	word = strings.TrimSpace(word)
	keywords := strings.FieldsFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-' && r != '_'
	})
	if len(keywords) == 0 {
		return db
	}

	whereSQLs := make([]string, 0, len(fields))
	whereArgs := make([]interface{}, 0, len(fields))
	for _, field := range fields {
		for _, kw := range keywords {
			// sqlite not support RLIKE
			// may have some special characters to ESCAPE
			var k = fmt.Sprintf("%s = ?", field)
			var v = strings.ReplaceAll(kw, "_", `\_`)
			if !exact {
				s := strings.ReplaceAll(strings.ToLower(kw), "_", `\_`)
				k = fmt.Sprintf("LOWER(%s) LIKE ?", field)
				v = "%" + s + "%"
			}
			whereSQLs = append(whereSQLs, k)
			whereArgs = append(whereArgs, v)
		}
	}
	return db.Where(strings.Join(whereSQLs, " OR "), whereArgs...)
}

func DBOrder(orders []Order, validOrderMap map[string]string) string {
	orderStrs := make([]string, 0, len(orders))
	for _, order := range orders {
		orderStr, valid := validOrderMap[order.Field]
		if valid {
			if order.Ascending {
				orderStr += " ASC"
			} else {
				orderStr += " DESC"
			}
			orderStrs = append(orderStrs, orderStr)
		}
	}
	return strings.Join(orderStrs, ", ")
}
